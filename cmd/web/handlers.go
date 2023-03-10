package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"subscription-service/database"
	"time"

	"github.com/phpdave11/gofpdf"                //This library allows us to create a pdf
	"github.com/phpdave11/gofpdf/contrib/gofpdi" //this library opens an existing pdf and uses it as a template
)

func (app *Config) homePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *Config) Loginpage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	// clean up session
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *Config) LandingPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "landing.page.gohtml", nil)
}

func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}
	//validate user
	u := database.User{
		Email:     r.Form.Get("email"),
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"),
		//specify the user is not an admin or not active
		Active:  0,
		IsAdmin: 0,
	}
	_, err = u.Insert(u)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to create a user")
		app.ErrorLog.Println("Unable to Create user")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}
	app.InfoLog.Println("User Registered")
	// send an activation email
	URL := fmt.Sprintf("http://localhost/activate?email=%s", u.Email)
	signedURL := GenerateTokenFromString(URL) //prevents url from being tampered
	app.InfoLog.Println(signedURL)
	app.InfoLog.Println("Activation emailUrl Generated")

	//create email message
	msg := Message{
		To:       u.Email,
		Subject:  "Activate your account by verifying mail",
		Template: "confirmation-email",
		Data:     template.HTML(signedURL), //cast the url into the html template
	}
	app.sendemail(msg)
	app.Session.Put(r.Context(), "flash", "Confirmation Mail Sent. Check your email to verify your mail.")
	app.InfoLog.Println("Activation Email Sent to Mail ")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

// send activation email to confirm we have the valid email address
func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	app.InfoLog.Println("Activate Email Page gottten")
	// validate url
	URL := r.RequestURI
	testURL := fmt.Sprintf("http://localhost%s", URL)
	okay := VerifyToken(testURL) //the url with the hash appended to it
	if !okay {
		app.Session.Put(r.Context(), "error", "invalid token")
		app.ErrorLog.Println("invalid token")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//activate account
	u, err := app.Models.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		app.Session.Put(r.Context(), "error", "No User Found")
		app.ErrorLog.Println("No User Found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u.Active = 1
	err = app.Models.User.Update(*u)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable User to update User.")
		app.ErrorLog.Println("Unable to update user")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.Session.Put(r.Context(), "flash", "Account activated. You can login.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	log.Println("Account Activated")
}

func (app *Config) PostLoginpage(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())

	// parse form post
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	// get email and password from form post
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, loginValue, _, err := app.Models.User.LoginUser(email, password)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid credentials.")
		app.ErrorLog.Println("Invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !loginValue {
		app.ErrorLog.Println("Invalid Password", password)
		msg := Message{
			To:      email,
			Subject: "Failed log in attempt",
			Data:    "Invalid Login ",
		}
		//call the channel func that sends email message easily
		app.sendemail(msg) //the app.Mailer.MailerChan which is a reciever from the message struct communictes this message sent to the
		//creating a session
		app.Session.Put(r.Context(), "error", "Invalid credentials.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)

	app.Session.Put(r.Context(), "flash", "Successful login!")
	app.InfoLog.Println("Succesful Login")
	http.Redirect(w, r, "/members/landing", http.StatusSeeOther)

}

var pathToManual = "./pdf"
var tempPath = "./temp"

func (app *Config) SubcribeToPlan(w http.ResponseWriter, r *http.Request) {
	// get the id of the plan that is chosen
	id := r.URL.Query().Get("id")

	planID, err := strconv.Atoi(id)
	if err != nil {
		app.ErrorLog.Println("Error Getting PlanID:", err)
		return
	}
	// get the plan from the database
	plan, err := app.Models.Plan.GetOne(planID)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to find plan")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
	}
	// get the user from the session
	user, ok := app.Session.Get(r.Context(), "user").(database.User)
	if !ok {
		app.Session.Put(r.Context(), "error", "Log in first")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// generate an invoice and email in it
	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done() //decrement waitgroup

		invoice, err := app.getInvoice(user, plan)
		if err != nil {
			app.ErrorChan <- err
		}
		app.InfoLog.Println("Invoice Generated")
		//send an email
		msg := Message{
			To:       user.Email,
			Subject:  "Your Invoice",
			Data:     invoice,
			Template: "invoice",
		}
		app.sendemail(msg)
		app.InfoLog.Println("Invoice Sent to Email")
	}()

	// generate a manual
	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()
		pdf := app.generateManual(user, plan)
		//write to a file
		err := pdf.OutputFileAndClose(fmt.Sprintf("%s/%d_manual.pdf", pathToManual, user.ID))
		if err != nil {
			app.ErrorChan <- err
			app.ErrorLog.Println("error in generating Manual", err)
			return
		}
		app.InfoLog.Println("Manual Generated")
		// send an email with the manual attached
		msg := Message{
			To:      user.Email,
			Subject: "Your Manual",
			Data:    "Your User Manual is attached",
			AttachmentMap: map[string]string{
				"Manual.pdf": fmt.Sprintf("%s/%d_manual.pdf", tempPath, user.ID),
			},
		}
		app.sendemail(msg)
		app.InfoLog.Println("Manual sent to email")
	}()
	// subscribe the user to an account
	err = app.Models.Plan.SubscribeUserToPlan(user, *plan)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Error Subscribing to Plan")
		app.ErrorLog.Println("Error in Subscribing to Plan")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}
	//update the user in the session from the db
	u, err := app.Models.User.GetOne(user.ID)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Error Getting User From Database")
		app.ErrorLog.Println("Error in getting User from db")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}
	app.Session.Put(r.Context(), "user", u)
	// redirect
	app.Session.Put(r.Context(), "flash", "Subscribed")
	http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
}

func (app *Config) getInvoice(u database.User, plan *database.Plan) (string, error) {
	app.InfoLog.Println("Amount is ", plan.PlanAmountFormatted)
	return plan.PlanAmountFormatted, nil //returns the price of the plan
}

func (app *Config) generateManual(u database.User, plan *database.Plan) *gofpdf.Fpdf {
	pdf := gofpdf.New("p", "mm", "Letter", "")
	pdf.SetMargins(10, 13, 10)

	importer := gofpdi.NewImporter()

	time.Sleep(5 * time.Second)
	t := importer.ImportPage(pdf, fmt.Sprintf("%s/manual.pdf", pathToManual), 1, "/MediaBox")
	pdf.AddPage() //we have a page already

	//use the imported template for the page
	importer.UseImportedTemplate(pdf, t, 0, 0, 215.9, 0)
	//set x and y coordinates
	pdf.SetX(75)
	pdf.SetY(150)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s %s", u.FirstName, u.LastName), "", "C", false)
	pdf.Ln(5) //performing line breaks

	//a cell that may span multiple linebreaks
	pdf.MultiCell(0, 4, fmt.Sprintf("%s User Guide", plan.PlanName), "", "C", false)

	return pdf
}

func (app *Config) ChoosePlans(w http.ResponseWriter, r *http.Request) {
	plans, err := app.Models.Plan.GetAll()
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	//pass the var to the template
	dataMap := make(map[string]any)
	dataMap["plans"] = plans
	app.render(w, r, "plans.page.gohtml", &TemplateData{
		Data: dataMap,
	})
}
