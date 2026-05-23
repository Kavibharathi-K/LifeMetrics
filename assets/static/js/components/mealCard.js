function renderMealCard(title, meal){

    const cards =
    document.querySelectorAll(".meal-card")

    cards.forEach(card => {

        const heading =
        card.querySelector("h3")

        if(heading.innerText !== title)
            return


        /*
        remove previous render
        */

        const oldList =
        card.querySelector(".meal-food-list")

        if(oldList) oldList.remove()


        const oldTotal =
        card.querySelector(".meal-total")

        if(oldTotal) oldTotal.remove()


        /*
        scrollable food container
        */

        const list =
        document.createElement("div")

        list.className =
        "meal-food-list"



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

            <div>Total</div>

            <div>

                ${Math.round(meal.totals.calories)} kcal

            </div>

        `



        /*
        insert before button
        */

        const button =
        card.querySelector("button")


        card.insertBefore(
            list,
            button
        )


        card.insertBefore(
            total,
            button
        )

    })

}