window.onload = function(){

    initializeTheme()

    const subPage =
    document.body.dataset.subpage

    console.log(
        "SubPage:",
        subPage
    )


    if(subPage === "today"){

        loadTodayNutrition()

        return

    }


    if(subPage === "goals"){

        loadGoals()

        const form =
        document.getElementById(
            "goals-form"
        )

        if(form){

            form.addEventListener(
                "submit",
                submitGoals
            )

        }

        return

    }


    if(subPage === "foods"){

        loadFoods()

    }


    if(subPage === "meals"){

        loadMeals()

    }


    if(subPage === "workout"){

        initializeWorkoutPage()

    }

}
