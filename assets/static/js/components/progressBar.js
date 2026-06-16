function updateProgress(
    id,
    value,
    goal,
    unit
){

    if(!goal) return

    const percent =
    (value / goal) * 100


    const bar =
    document.getElementById(
        `${id}-progress`
    )

    const text =
    document.getElementById(
        `${id}-remaining`
    )


    /* width */

    bar.style.width =
    Math.min(percent, 100) + "%"



    /* status text */

    if(value < goal){

        const remaining =
        goal - value

        text.innerText =
        remaining.toFixed(0) +
        unit +
        " left"

        bar.style.background =
        "var(--accent)"

    }

    else if(value === goal){

        text.innerText =
        "Goal reached"

        bar.style.background =
        "#16a34a"

    }

    else{

        const extra =
        value - goal

        text.innerText =
        extra.toFixed(0) +
        unit +
        " extra"

        bar.style.background =
        "#dc2626"


        bar.style.width =
        Math.min(percent, 130) + "%"

    }

}