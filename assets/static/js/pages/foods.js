async function loadFoods(){

    try{

        const foods =
        await getFoods()

        console.log("foods:", foods)

        const tableBody =
        document.querySelector(
            "#foodTable tbody"
        )

        tableBody.innerHTML = ""

        foods.forEach(food => {

            const row = `

                <tr>

                    <td>${food.name}</td>

                    <td>
                        ${food.base_quantity}
                        ${food.measurement_type}
                    </td>

                    <td>${food.calories || 0}</td>

                    <td>${food.protein || 0}</td>

                    <td>${food.carbs || 0}</td>

                    <td>${food.fat || 0}</td>

                </tr>

            `

            tableBody.innerHTML += row

        })

    }
    catch(error){

        console.error(
            "Error loading foods:",
            error
        )

    }

}


function showFoodsView(){

    document.getElementById(
        "pageTitle"
    ).innerText =
    "Available Foods"

    document.getElementById(
        "foodsView"
    ).style.display =
    "block"

    document.getElementById(
        "mealsView"
    ).style.display =
    "none"

    loadFoods()

}