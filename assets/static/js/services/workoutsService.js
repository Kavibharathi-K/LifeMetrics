async function getWorkoutSchedules(){

    return apiGet(
        "/workouts/schedules"
    )

}


async function getWorkoutEmailSettings(){

    return apiGet(
        "/workouts/email-settings"
    )

}


async function updateWorkoutEmailSettings(
    data
){

    return apiPut(
        "/workouts/email-settings",
        data
    )

}


async function createWorkoutSchedule(
    data
){

    return apiPost(
        "/workouts/schedules",
        data
    )

}


async function updateWorkoutSchedule(
    workoutScheduleId,
    data
){

    return apiPut(
        `/workouts/schedules/${workoutScheduleId}`,
        data
    )

}


async function addWorkoutScheduleExercise(
    workoutScheduleId,
    data
){

    return apiPost(
        `/workouts/schedules/${workoutScheduleId}/exercises`,
        data
    )

}


async function deleteWorkoutSchedule(
    workoutScheduleId
){

    return apiDelete(
        `/workouts/schedules/${workoutScheduleId}`
    )

}