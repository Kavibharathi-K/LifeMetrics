function initializeTheme(){

    const savedTheme =
    localStorage.getItem("theme")

    if(savedTheme === "dark"){

        document.body.classList.add(
            "dark-mode"
        )

    }


    const button =
    document.getElementById(
        "themeToggle"
    )

    if(!button) return


    updateThemeButtonText()


    button.addEventListener(
        "click",
        toggleTheme
    )

}


function toggleTheme(){

    document.body.classList.toggle(
        "dark-mode"
    )

    const isDark =
    document.body.classList.contains(
        "dark-mode"
    )

    localStorage.setItem(
        "theme",
        isDark ? "dark" : "light"
    )

    updateThemeButtonText()

}


function updateThemeButtonText(){

    const button =
    document.getElementById(
        "themeToggle"
    )

    if(!button) return


    const isDark =
    document.body.classList.contains(
        "dark-mode"
    )

    button.innerText =
    isDark
        ? "☀️ Light"
        : "🌙 Dark"

}