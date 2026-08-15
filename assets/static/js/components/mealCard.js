function renderMealCard(title, meal){

    const cards =
    document.querySelectorAll(".meal-card")

    cards.forEach(card => {

        const heading =
        card.querySelector("h2")

        if(heading.innerText !== title)
            return

        const mealContainer = card.querySelector(".meal-foods")
        if(!mealContainer) return


        /*
        remove previous render
        */

        mealContainer.innerHTML = ""


        /*
        scrollable food container
        */

        const list =
        document.createElement("div")

        list.className =
        "meal-food-list"


        const hero =
        document.createElement("div")

        hero.className =
        "meal-hero"

        hero.innerHTML = `

            <div class="meal-hero-value">

                ${Math.round(meal.totals.calories)}

            </div>

            <div class="meal-hero-label">

                kcal today

            </div>

        `

        /*
        if no foods
        */

        if(!meal.foods || meal.foods.length === 0){

            list.innerHTML =

            `<div class="meal-empty">
                No foods added yet
            </div>`

        }


        /*
        render foods
        */

        else{

            meal.foods.forEach(food => {

                const row =
                document.createElement("div")

                row.className =
                "meal-food-row"


                row.innerHTML = `

                    <div>

                        <div class="meal-food-name">

                            ${food.name}

                        </div>

                        <div class="meal-food-qty">

                            ${food.quantity} ${food.unit}

                        </div>

                    </div>


                    <div class="meal-food-cal">

                        ${Math.round(food.calories)} kcal

                    </div>

                `


                list.appendChild(row)

            })

        }



        /*
        total section
        */

        const total =
        document.createElement("div")

        total.className =
        "meal-total"

        total.innerHTML = `

            <div>Total Foods</div>

            <div>

                ${meal.foods ? meal.foods.length : 0}

            </div>

        `



       mealContainer.appendChild(hero)
       mealContainer.appendChild(list)
       mealContainer.appendChild(total)

    })

}