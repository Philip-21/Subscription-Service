(function () {
      'use strict' //enable strict mode

      // method call selects all elements on the page with a needs-validation class and returns them as a NodeList.
      let forms = document.querySelectorAll('.needs-validation')

      //converts the NodeList to an array so that the forEach method can be used to iterate over it
      Array.prototype.slice.call(forms)
          .forEach(function (form) {
              form.addEventListener('submit', function (event) {
                  if (!form.checkValidity()) {
                      event.preventDefault()
                      event.stopPropagation()
                  }

                  form.classList.add('was-validated')
              }, false)
          })
  })()