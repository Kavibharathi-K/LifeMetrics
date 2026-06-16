const workoutDays = [
    "Monday",
    "Tuesday",
    "Wednesday",
    "Thursday",
    "Friday",
    "Saturday",
    "Sunday"
]

let workoutSchedules = []
let workoutPageInitialized = false

function initializeWorkoutPage() {

    if(workoutPageInitialized){
        return
    }

    const grid =
    document.getElementById(
        "workoutScheduleGrid"
    )

    if(!grid){
        return
    }

    workoutPageInitialized = true

    document.getElementById(
        "openWorkoutModalBtn"
    )
    .addEventListener(
        "click",
        () => openWorkoutModal()
    )

    document.getElementById(
        "addExerciseInputBtn"
    )
    .addEventListener(
        "click",
        () => addWorkoutExerciseInput()
    )

    document.getElementById(
        "workoutForm"
    )
    .addEventListener(
        "submit",
        submitWorkoutForm
    )

    document.getElementById(
        "addWorkoutExerciseForm"
    )
    .addEventListener(
        "submit",
        submitWorkoutExercise
    )

    document
    .querySelectorAll(
        "[data-close-workout-modal]"
    )
    .forEach(button => {
        button.addEventListener(
            "click",
            closeWorkoutModal
        )
    })

    document
    .querySelectorAll(
        "[data-close-exercise-modal]"
    )
    .forEach(button => {
        button.addEventListener(
            "click",
            closeWorkoutExerciseModal
        )
    })

    document
    .querySelectorAll(
        ".workout-modal"
    )
    .forEach(modal => {
        modal.addEventListener(
            "click",
            event => {
                if(event.target === modal){
                    closeWorkoutModal()
                    closeWorkoutExerciseModal()
                }
            }
        )
    })

    document.addEventListener(
        "keydown",
        event => {
            if(event.key === "Escape"){
                closeWorkoutModal()
                closeWorkoutExerciseModal()
            }
        }
    )

    loadWorkoutSchedules()
}

if(document.readyState === "loading"){

    document.addEventListener(
        "DOMContentLoaded",
        initializeWorkoutPage
    )

}
else{

    initializeWorkoutPage()

}

async function loadWorkoutSchedules() {

    setWorkoutPageMessage("")

    try{

        workoutSchedules =
        await getWorkoutSchedules()

        renderWorkoutSchedules()

    }
    catch(error){

        document.getElementById(
            "workoutScheduleGrid"
        ).innerHTML = `
            <div class="workout-loading workout-load-error">
                Could not load your workout schedule.
                <button type="button" onclick="loadWorkoutSchedules()">
                    Try Again
                </button>
            </div>
        `

        setWorkoutPageMessage(
            error.message,
            "error"
        )

    }
}

function renderWorkoutSchedules() {

    const grid =
    document.getElementById(
        "workoutScheduleGrid"
    )

    grid.innerHTML = ""

    workoutDays.forEach(
        (dayName, index) => {

            const dayOfWeek =
            index + 1

            const schedule =
            workoutSchedules.find(
                item =>
                item.day_of_week === dayOfWeek
            )

            grid.appendChild(
                createWorkoutDayCard(
                    dayName,
                    dayOfWeek,
                    schedule
                )
            )

        }
    )
}

function createWorkoutDayCard(
    dayName,
    dayOfWeek,
    schedule
) {

    const card =
    document.createElement("article")

    card.className =
    schedule
        ? "workout-day-card has-workout"
        : "workout-day-card rest-day"

    const header =
    document.createElement("div")

    header.className =
    "workout-day-card-header"

    const dayBlock =
    document.createElement("div")

    const dayLabel =
    document.createElement("span")

    dayLabel.className =
    "workout-day-label"

    dayLabel.textContent =
    `Day ${dayOfWeek}`

    const heading =
    document.createElement("h2")

    heading.textContent =
    dayName

    dayBlock.append(
        dayLabel,
        heading
    )

    const badge =
    document.createElement("span")

    badge.className =
    "workout-status-badge"

    badge.textContent =
    schedule
        ? `${schedule.exercises.length} exercises`
        : "Rest day"

    header.append(
        dayBlock,
        badge
    )

    card.appendChild(header)

    if(!schedule){

        const emptyState =
        document.createElement("div")

        emptyState.className =
        "workout-empty-state"

        const text =
        document.createElement("p")

        text.textContent =
        "No workout planned. Keep it as recovery or add a session."

        const button =
        document.createElement("button")

        button.type = "button"
        button.className =
        "workout-secondary-button"

        button.textContent =
        "+ Add Workout"

        button.addEventListener(
            "click",
            () => openWorkoutModal(dayOfWeek)
        )

        emptyState.append(
            text,
            button
        )

        card.appendChild(emptyState)

        return card
    }

    const workoutName =
    document.createElement("h3")

    workoutName.className =
    "workout-name"

    workoutName.textContent =
    schedule.workout_name

    card.appendChild(workoutName)

    const exerciseList =
    document.createElement("ol")

    exerciseList.className =
    "workout-exercise-list"

    if(schedule.exercises.length === 0){

        const emptyExercise =
        document.createElement("li")

        emptyExercise.className =
        "workout-no-exercises"

        emptyExercise.textContent =
        "No exercises added yet."

        exerciseList.appendChild(
            emptyExercise
        )

    }
    else{

        schedule.exercises.forEach(
            exercise => {

                const item =
                document.createElement("li")

                const order =
                document.createElement("span")

                order.className =
                "workout-exercise-order"

                order.textContent =
                exercise.exercise_order

                const name =
                document.createElement("span")

                name.textContent =
                exercise.exercise_name

                item.append(
                    order,
                    name
                )

                exerciseList.appendChild(
                    item
                )
            }
        )

    }

    card.appendChild(exerciseList)

    const actions =
    document.createElement("div")

    actions.className =
    "workout-card-actions"

    actions.append(
        createWorkoutActionButton(
            "+ Exercise",
            "workout-secondary-button",
            () => openWorkoutExerciseModal(
                schedule
            )
        ),
        createWorkoutActionButton(
            "Edit",
            "workout-secondary-button",
            () => openWorkoutModal(
                schedule.day_of_week,
                schedule
            )
        ),
        createWorkoutActionButton(
            "Delete",
            "workout-danger-button",
            () => confirmDeleteWorkout(
                schedule
            )
        )
    )

    card.appendChild(actions)

    return card
}

function createWorkoutActionButton(
    label,
    className,
    onClick
) {

    const button =
    document.createElement("button")

    button.type = "button"
    button.className = className
    button.textContent = label
    button.addEventListener(
        "click",
        onClick
    )

    return button
}

function openWorkoutModal(
    dayOfWeek = null,
    schedule = null
) {

    const modal =
    document.getElementById(
        "workoutModal"
    )

    const isEditing =
    Boolean(schedule)

    document.getElementById(
        "workoutModalTitle"
    ).textContent =
    isEditing
        ? "Edit Workout"
        : "Add Workout"

    document.getElementById(
        "saveWorkoutBtn"
    ).textContent =
    isEditing
        ? "Save Changes"
        : "Save Workout"

    document.getElementById(
        "workoutScheduleId"
    ).value =
    schedule
        ? schedule.workout_schedule_id
        : ""

    document.getElementById(
        "workoutNameInput"
    ).value =
    schedule
        ? schedule.workout_name
        : ""

    populateWorkoutDayOptions(
        dayOfWeek,
        schedule
    )

    const exerciseContainer =
    document.getElementById(
        "workoutExerciseInputs"
    )

    exerciseContainer.innerHTML = ""

    const exercises =
    schedule && schedule.exercises.length
        ? schedule.exercises.map(
            exercise =>
            exercise.exercise_name
        )
        : [""]

    exercises.forEach(
        exerciseName =>
        addWorkoutExerciseInput(
            exerciseName
        )
    )

    setWorkoutFormError("")
    setWorkoutButtonLoading(false)

    modal.style.display = "flex"
    modal.setAttribute(
        "aria-hidden",
        "false"
    )

    setTimeout(
        () => {
            document.getElementById(
                "workoutNameInput"
            ).focus()
        },
        0
    )
}

function closeWorkoutModal() {

    const modal =
    document.getElementById(
        "workoutModal"
    )

    if(!modal){
        return
    }

    modal.style.display = "none"
    modal.setAttribute(
        "aria-hidden",
        "true"
    )
}

function populateWorkoutDayOptions(
    selectedDay,
    editingSchedule
) {

    const select =
    document.getElementById(
        "workoutDayInput"
    )

    select.innerHTML = ""

    const firstAvailableDay =
    workoutDays.findIndex(
        (_, index) =>
        !workoutSchedules.some(
            schedule =>
            schedule.day_of_week ===
            index + 1
        )
    ) + 1

    const activeDay =
    selectedDay ||
    (firstAvailableDay > 0
        ? firstAvailableDay
        : 1)

    workoutDays.forEach(
        (dayName, index) => {

            const dayOfWeek =
            index + 1

            const occupied =
            workoutSchedules.some(
                schedule =>
                schedule.day_of_week ===
                dayOfWeek &&
                (!editingSchedule ||
                 schedule.workout_schedule_id !==
                 editingSchedule.workout_schedule_id)
            )

            const option =
            document.createElement("option")

            option.value = dayOfWeek
            option.textContent =
            occupied
                ? `${dayName} - already planned`
                : dayName

            option.disabled = occupied
            option.selected =
            dayOfWeek === activeDay

            select.appendChild(option)
        }
    )
}

function addWorkoutExerciseInput(
    value = ""
) {

    const container =
    document.getElementById(
        "workoutExerciseInputs"
    )

    const row =
    document.createElement("div")

    row.className =
    "workout-exercise-input-row"

    const order =
    document.createElement("span")

    order.className =
    "workout-input-order"

    const input =
    document.createElement("input")

    input.type = "text"
    input.maxLength = 150
    input.placeholder =
    "Exercise name"
    input.value = value
    input.className =
    "workout-exercise-name-input"

    const moveUp =
    createWorkoutInputButton(
        "Up",
        "Move exercise up",
        () => moveWorkoutExerciseInput(
            row,
            -1
        )
    )

    const moveDown =
    createWorkoutInputButton(
        "Down",
        "Move exercise down",
        () => moveWorkoutExerciseInput(
            row,
            1
        )
    )

    const remove =
    createWorkoutInputButton(
        "Remove",
        "Remove exercise",
        () => {
            row.remove()
            ensureWorkoutExerciseInput()
            updateWorkoutExerciseInputOrder()
        }
    )

    remove.classList.add(
        "workout-remove-input"
    )

    row.append(
        order,
        input,
        moveUp,
        moveDown,
        remove
    )

    container.appendChild(row)
    updateWorkoutExerciseInputOrder()

    if(value === ""){
        input.focus()
    }
}

function createWorkoutInputButton(
    label,
    ariaLabel,
    onClick
) {

    const button =
    document.createElement("button")

    button.type = "button"
    button.className =
    "workout-input-action"
    button.textContent = label
    button.setAttribute(
        "aria-label",
        ariaLabel
    )
    button.addEventListener(
        "click",
        onClick
    )

    return button
}

function moveWorkoutExerciseInput(
    row,
    direction
) {

    const sibling =
    direction < 0
        ? row.previousElementSibling
        : row.nextElementSibling

    if(!sibling){
        return
    }

    if(direction < 0){
        row.parentElement.insertBefore(
            row,
            sibling
        )
    }
    else{
        row.parentElement.insertBefore(
            sibling,
            row
        )
    }

    updateWorkoutExerciseInputOrder()
}

function ensureWorkoutExerciseInput() {

    const container =
    document.getElementById(
        "workoutExerciseInputs"
    )

    if(container.children.length === 0){
        addWorkoutExerciseInput()
    }
}

function updateWorkoutExerciseInputOrder() {

    document
    .querySelectorAll(
        ".workout-exercise-input-row"
    )
    .forEach(
        (row, index) => {
            row.querySelector(
                ".workout-input-order"
            ).textContent =
            index + 1
        }
    )
}

async function submitWorkoutForm(event) {

    event.preventDefault()

    const scheduleId =
    document.getElementById(
        "workoutScheduleId"
    ).value

    const workoutName =
    document.getElementById(
        "workoutNameInput"
    ).value.trim()

    const exercises =
    Array.from(
        document.querySelectorAll(
            ".workout-exercise-name-input"
        )
    )
    .map(input => input.value.trim())
    .filter(Boolean)

    if(!workoutName){
        setWorkoutFormError(
            "Enter a workout name."
        )
        return
    }

    if(exercises.length === 0){
        setWorkoutFormError(
            "Add at least one exercise."
        )
        return
    }

    const data = {
        day_of_week:
        Number(
            document.getElementById(
                "workoutDayInput"
            ).value
        ),
        workout_name:
        workoutName,
        exercises:
        exercises
    }

    setWorkoutFormError("")
    setWorkoutButtonLoading(true)

    try{

        if(scheduleId){
            await updateWorkoutSchedule(
                scheduleId,
                data
            )
        }
        else{
            await createWorkoutSchedule(
                data
            )
        }

        closeWorkoutModal()

        setWorkoutPageMessage(
            scheduleId
                ? "Workout updated."
                : "Workout added.",
            "success"
        )

        await loadWorkoutSchedules()

    }
    catch(error){

        setWorkoutFormError(
            error.message
        )

    }
    finally{

        setWorkoutButtonLoading(false)

    }
}

function openWorkoutExerciseModal(
    schedule
) {

    document.getElementById(
        "exerciseWorkoutScheduleId"
    ).value =
    schedule.workout_schedule_id

    document.getElementById(
        "addWorkoutExerciseTitle"
    ).textContent =
    `Add to ${schedule.workout_name}`

    document.getElementById(
        "newWorkoutExerciseName"
    ).value = ""

    document.getElementById(
        "addWorkoutExerciseError"
    ).textContent = ""

    const modal =
    document.getElementById(
        "addWorkoutExerciseModal"
    )

    modal.style.display = "flex"
    modal.setAttribute(
        "aria-hidden",
        "false"
    )

    setTimeout(
        () => {
            document.getElementById(
                "newWorkoutExerciseName"
            ).focus()
        },
        0
    )
}

function closeWorkoutExerciseModal() {

    const modal =
    document.getElementById(
        "addWorkoutExerciseModal"
    )

    if(!modal){
        return
    }

    modal.style.display = "none"
    modal.setAttribute(
        "aria-hidden",
        "true"
    )
}

async function submitWorkoutExercise(
    event
) {

    event.preventDefault()

    const scheduleId =
    document.getElementById(
        "exerciseWorkoutScheduleId"
    ).value

    const exerciseName =
    document.getElementById(
        "newWorkoutExerciseName"
    ).value.trim()

    const errorElement =
    document.getElementById(
        "addWorkoutExerciseError"
    )

    if(!exerciseName){
        errorElement.textContent =
        "Enter an exercise name."
        return
    }

    const button =
    document.getElementById(
        "saveWorkoutExerciseBtn"
    )

    button.disabled = true
    button.textContent = "Adding..."
    errorElement.textContent = ""

    try{

        await addWorkoutScheduleExercise(
            scheduleId,
            {
                exercise_name:
                exerciseName
            }
        )

        closeWorkoutExerciseModal()

        setWorkoutPageMessage(
            "Exercise added.",
            "success"
        )

        await loadWorkoutSchedules()

    }
    catch(error){

        errorElement.textContent =
        error.message

    }
    finally{

        button.disabled = false
        button.textContent =
        "Add Exercise"

    }
}

async function confirmDeleteWorkout(
    schedule
) {

    const confirmed =
    window.confirm(
        `Delete ${workoutDays[schedule.day_of_week - 1]}'s ${schedule.workout_name} workout?`
    )

    if(!confirmed){
        return
    }

    try{

        await deleteWorkoutSchedule(
            schedule.workout_schedule_id
        )

        setWorkoutPageMessage(
            "Workout deleted.",
            "success"
        )

        await loadWorkoutSchedules()

    }
    catch(error){

        setWorkoutPageMessage(
            error.message,
            "error"
        )

    }
}

function setWorkoutButtonLoading(
    loading
) {

    const button =
    document.getElementById(
        "saveWorkoutBtn"
    )

    if(!button){
        return
    }

    button.disabled = loading

    if(loading){
        button.dataset.label =
        button.textContent
        button.textContent =
        "Saving..."
    }
    else if(button.dataset.label){
        button.textContent =
        button.dataset.label
        delete button.dataset.label
    }
}

function setWorkoutFormError(
    message
) {

    document.getElementById(
        "workoutFormError"
    ).textContent =
    message
}

function setWorkoutPageMessage(
    message,
    type = ""
) {

    const element =
    document.getElementById(
        "workoutPageMessage"
    )

    if(!element){
        return
    }

    element.textContent = message
    element.className =
    "page-message workout-page-message"

    if(type){
        element.classList.add(type)
    }
}
