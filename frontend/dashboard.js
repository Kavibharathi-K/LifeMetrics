function toggleMenu(id){

    const menu = document.getElementById(id)

    if(menu.style.display === "block"){
        menu.style.display = "none"
    }else{
        menu.style.display = "block"
    }
}


function loadPage(page){

const frame =
document.getElementById("contentFrame")

/* remove active from all submenu items */

document
.querySelectorAll(".submenu-item")
.forEach(item => {

item.classList.remove("active")

})


/* set active on selected menu */

event.target.classList.add("active")



if(page === "add-food"){

frame.src = "/add_food.html"

}


if(page === "food-list"){

frame.src = "/food_list.html"

}

}