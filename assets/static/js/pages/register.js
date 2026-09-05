document.addEventListener(
    "DOMContentLoaded",
    function(){

        if(isAuthenticated()){

            window.location.href = "/"

            return

        }

        const registerForm =
        document.getElementById(
            "registerForm"
        )

        const registerError =
        document.getElementById(
            "registerError"
        )

        const registerButton =
        document.getElementById(
            "registerButton"
        )


        registerForm.addEventListener(
            "submit",
            async function(event){

                event.preventDefault()

                registerError.innerText = ""


                const email =
                document.getElementById(
                    "registerEmail"
                ).value.trim()

                const password =
                document.getElementById(
                    "registerPassword"
                ).value


                registerButton.disabled = true

                registerButton.innerText =
                "Creating account..."


                try{

                    await register({
                        email,
                        password
                    })


                    window.location.href =
                    "/login"

                }
                catch(error){

                    console.error(
                        "Registration error:",
                        error
                    )

                    registerError.innerText =
                    error.message ||
                    "Unable to create account"

                }
                finally{

                    registerButton.disabled = false

                    registerButton.innerText =
                    "Create account"

                }

            }
        )

    }
)