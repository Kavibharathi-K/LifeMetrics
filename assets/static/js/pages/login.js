document.addEventListener(
    "DOMContentLoaded",
    function(){

        if(isAuthenticated()){
            window.location.href = "/"
            return
        }

        const loginForm =
        document.getElementById(
            "loginForm"
        )

        const loginError =
        document.getElementById(
            "loginError"
        )

        const loginButton =
        document.getElementById(
            "loginButton"
        )

        loginForm.addEventListener(
            "submit",
            async function(event){

                event.preventDefault()
                loginError.innerText = ""

                const email =
                document.getElementById(
                    "loginEmail"
                ).value.trim()

                const password =
                document.getElementById(
                    "loginPassword"
                ).value

                loginButton.disabled = true
                loginButton.innerText = "Logging in..."

                try{
                    await login({email,password})
                    window.location.href = "/"
                }
                catch(error){
                    console.error("Login error:", error)

                    loginError.innerText =
                    error.message ||
                    "Invalid email or password"
                }
                finally{
                    loginButton.disabled = false
                    loginButton.innerText = "Login"
                }
            }
        )
    }
)