let selectedFood = null
let searchResults = []
let debounceTimer = null

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
                    <td>${food.fiber || 0}</td>

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

function openAddFoodModal() {

    const modal =
    document.getElementById(
        "addFoodModal"
    )

    modal.style.display = "flex"
    modal.setAttribute(
        "aria-hidden",
        "false"
    )

    const input =
    document.getElementById(
        "foodSearchInput"
    )

    input.focus()

}

function closeAddFoodModal() {

    const modal =
    document.getElementById(
        "addFoodModal"
    )

    modal.style.display = "none"
    modal.setAttribute(
        "aria-hidden",
        "true"
    )

    document.getElementById(
        "foodSearchInput"
    ).value = ""

    document.getElementById(
        "foodSearchResults"
    ).innerHTML = ""

    selectedFood = null

    document.getElementById(
        "submitFoodBtn"
    ).disabled = true

}

function initializeFoodSearch() {

    const input =
    document.getElementById(
        "foodSearchInput"
    )

    if(!input){
        return
    }

    input.addEventListener(
        "input",
        function(){

            clearTimeout(
                debounceTimer
            )

            debounceTimer =
            setTimeout(
                async () => {

                    const query =
                    input.value.trim()

                    if(query.length < 2){

                        document.getElementById(
                            "foodSearchResults"
                        ).innerHTML = ""

                        return

                    }

                    await searchFoods(
                        query
                    )

                },
                400
            )

        }
    )

}

async function searchFoods(query){

    try{

        searchResults =
        await searchUSDAFoods(
            query
        )

        renderSearchResults()

    }
    catch(error){

        console.error(
            error
        )

    }

}

function renderSearchResults(){

    const container =
    document.getElementById(
        "foodSearchResults"
    )

    container.innerHTML = ""

    searchResults.forEach(
        (
            food,
            index
        ) => {

            const div =
            document.createElement(
                "div"
            )

            div.className =
            "food-search-item"

            div.innerHTML = `

                <strong>
                    ${food.name}
                </strong>

                <div>

                    ${food.calories} kcal |

                    ${food.protein} P |

                    ${food.carbs} C |

                    ${food.fat} F

                </div>

            `

            div.onclick =
            () => selectFood(
                index
            )

            container.appendChild(
                div
            )

        }
    )

}

function selectFood(index){

    selectedFood =
    searchResults[index]

    document
    .querySelectorAll(
        ".food-search-item"
    )
    .forEach(
        item =>
        item.classList.remove(
            "selected"
        )
    )

    document
    .querySelectorAll(
        ".food-search-item"
    )
    [index]
    .classList.add(
        "selected"
    )

    document.getElementById(
        "submitFoodBtn"
    ).disabled = false

}

// async function submitSelectedFood(){

//     console.log(
//         selectedFood
//     )

// }

async function submitSelectedFood() {

    console.log("selectedFood", selectedFood);

    const response = await fetch("/foods/addfood", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            name: selectedFood.name,
            measurement_type: "gram",
            base_quantity: selectedFood.base_quantity,
            calories: selectedFood.calories,
            protein: selectedFood.protein,
            carbs: selectedFood.carbs,
            fat: selectedFood.fat,
            fiber: selectedFood.fiber
        })
    });

    const data = await response.json();

    console.log(data);
}

document.addEventListener(
    "DOMContentLoaded",
    function(){

        initializeFoodSearch()

    }
)