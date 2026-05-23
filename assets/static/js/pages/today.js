async function loadTodayNutrition(){

    try{

        await loadMacroGoals()

        const data =
        await getTodayNutrition()

        const nutrition =
        data[0]

        if(!nutrition) return


        const calories =
        Math.round(
            nutrition.total_calories
        )

        const protein =
        nutrition.total_protein

        const carbs =
        nutrition.total_carbs

        const fat =
        nutrition.total_fat

        const fiber =
        nutrition.total_fiber


        document
        .getElementById(
            "calories-value"
        )
        .innerText =
        calories


        document
        .getElementById(
            "protein-value"
        )
        .innerText =
        protein.toFixed(1) + " g"


        document
        .getElementById(
            "carbs-value"
        )
        .innerText =
        carbs.toFixed(1) + " g"


        document
        .getElementById(
            "fat-value"
        )
        .innerText =
        fat.toFixed(1) + " g"


        document
        .getElementById(
            "fiber-value"
        )
        .innerText =
        fiber.toFixed(1) + " g"


        updateProgress(
            "calories",
            calories,
            goalsData.calories,
            " kcal"
        )

        updateProgress(
            "protein",
            protein,
            goalsData.protein,
            " g"
        )

        updateProgress(
            "carbs",
            carbs,
            goalsData.carbs,
            " g"
        )

        updateProgress(
            "fat",
            fat,
            goalsData.fat,
            " g"
        )

        updateProgress(
            "fiber",
            fiber,
            goalsData.fiber,
            " g"
        )

    }
    catch(err){

        console.error(
            "Failed dashboard load",
            err
        )

    }

}