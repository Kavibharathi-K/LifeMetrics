async function workoutRequest(url, options = {}) {

    const response =
    await fetch(url, options)

    if(response.status === 204){
        return null
    }

    const data =
    await response.json()
    .catch(() => ({}))

    if(!response.ok){
        throw new Error(
            data.error ||
            "Workout request failed"
        )
    }

    return data
}

async function getWorkoutSchedules() {

    return workoutRequest(
        "/workouts/schedules"
    )
}

async function createWorkoutSchedule(data) {

    return workoutRequest(
        "/workouts/schedules",
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        }
    )
}

async function updateWorkoutSchedule(
    workoutScheduleId,
    data
) {

    return workoutRequest(
        `/workouts/schedules/${workoutScheduleId}`,
        {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        }
    )
}

async function addWorkoutScheduleExercise(
    workoutScheduleId,
    data
) {

    return workoutRequest(
        `/workouts/schedules/${workoutScheduleId}/exercises`,
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        }
    )
}

async function deleteWorkoutSchedule(
    workoutScheduleId
) {

    return workoutRequest(
        `/workouts/schedules/${workoutScheduleId}`,
        {
            method: "DELETE"
        }
    )
}
