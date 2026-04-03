const API_URL = "http://localhost:8080/foods"


/* ADD FOOD */
const form =
document.getElementById("foodForm")

if(form){

form.addEventListener("submit", async (e) => {

e.preventDefault()

const payload = {

name:
document.getElementById("name").value,

measurement_type:
document.getElementById("measurement_type").value,

base_quantity:
Number(document.getElementById("base_quantity").value),

calories:
Number(document.getElementById("calories").value),

protein:
Number(document.getElementById("protein").value),

carbs:
Number(document.getElementById("carbs").value),

fat:
Number(document.getElementById("fat").value),

fiber:
Number(document.getElementById("fiber").value)

}

await fetch(`${API_URL}/addfood`,{

method:"POST",

headers:{
"Content-Type":"application/json"
},

body:JSON.stringify(payload)

})

alert("Food Added")

form.reset()

})

}



/* FOOD LIST */
const table =
document.getElementById("foodTable")

if(table){

loadFoods()

}


async function loadFoods(){

const res =
await fetch(`${API_URL}/listfoods`)

const foods =
await res.json()

table.innerHTML = ""

foods.forEach(food => {

table.innerHTML += `

<tr>

<td>${food.name}</td>

<td>
${food.base_quantity}
${food.measurement_type}
</td>

<td>
${food.calories}
</td>

<td>

<button onclick="editFood(${food.food_id})">

Edit

</button>

</td>

</tr>

`

})

}


function editFood(id){

window.location.href =
`add_food.html?food_id=${id}`

}