async function getLatestGoals(){

    return apiGet(
        "/usermetrics/getlatestusermetrics"
    )

}


async function saveGoals(data){

    return apiPost(
        "/usermetrics/addusermetrics",
        data
    )

}