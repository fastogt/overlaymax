function updatePage(domain, pid) {
    let list = document.getElementById("show_logos");
    let btnUpdate = document.getElementById("btnUpdate")
    let received_msg = new Object();
    btnUpdate.textContent = "Apply"
    received_msg.players = [{
        "team": document.getElementById("player_id_00").value,
        "score": parseInt(document.getElementById("player_id_score_00").value),
        "logo": document.getElementById("player_logo_0").src
    },
    {
        "team": document.getElementById("player_id_11").value,
        "score": parseInt(document.getElementById("player_id_score_11").value),
        "logo": document.getElementById("player_logo_1").src
    }];
    let local_time = new Date().getTime()
    let stadium = document.getElementById("local_stadium")
    received_msg.date_time_location = {
        "local_time": local_time,
        "local_stadium": stadium.textContent // FIXME
    }
    received_msg.bg_color = "green"
    received_msg.id = pid
    if (list.options[list.selectedIndex].value == 1) {
        received_msg.show_logos = false
    } else {
        received_msg.show_logos = true 
    }
    $.ajax({
        type: 'POST',
        url: domain + "/overlay/football/create",
        data: JSON.stringify(received_msg),
        contentType: "application/json",
        dataType: 'json',
        success: function () {}
    })
}

