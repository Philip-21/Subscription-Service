package main

//a  helper wrapper to send email easily
func (app *Config) sendemail(msg Message) {
	//add counter to waitgroup , increment wg by 1
	app.Wait.Add(1)
	app.Mailer.MailerChan <- msg
}
