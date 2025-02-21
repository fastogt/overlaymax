function convertLocalDate(localTime) {
    let convert_local_time = new Date(parseInt(localTime))
    let local_date = (convert_local_time.getDate()) + '.'
        + ('0' + (convert_local_time.getMonth() + 1)).slice(-2) + '.'
        + convert_local_time.getFullYear();
    let local_time = ('0' + (convert_local_time.getHours())).slice(-2) + ':'
        + ('0' + (convert_local_time.getMinutes())).slice(-2);
    document.getElementById("local_date").textContent = local_date;
    document.getElementById("local_time").textContent = local_time;
}


function changeOverlay(received_msg) {
    print(received_msg)
    // Players name
    document.getElementById("player_id_0").textContent = received_msg["players"][0]["team"];
    document.getElementById("player_id_1").textContent = received_msg["players"][1]["team"];
    //Score
    document.getElementById("player_id_score_0").textContent = received_msg["players"][0]["score"];
    document.getElementById("player_id_score_1").textContent = received_msg["players"][1]["score"];
    //Logo
    document.getElementById("player_logo_0").src = received_msg["players"][0]["logo"];
    document.getElementById("player_logo_1").src = received_msg["players"][1]["logo"];
    // DATE TIME STADIUM
    var dt = received_msg["date_time_location"];
    convertLocalDate(dt["local_time"].toString());
    document.getElementById("local_stadium").textContent = dt["local_stadium"];

    // BG_COLOR
    document.getElementById("downline_up").style.backgroundColor = received_msg["bg_color"]
    document.body.style.backgroundColor = received_msg["bg_color"]
    document.getElementById("btnUpdate").textContent = received_msg["btn_name"]
    // Show, hide logos
    if (received_msg.show_logos == false) {
       document.getElementById("total_info").style.display = 'none';
    } else {
       document.getElementById("total_info").style.display = 'block';
    }
}


let socketWS = null

function startSocket(overlay){
    let stabled = overlay.replace('http', 'ws')
    socketWS = new WebSocket(stabled)
    socketWS.onopen = function () {};

    socketWS.onmessage = async function (event) {
        let received_msg = JSON.parse(event.data)
        if (received_msg["type"] === "update") {
            changeOverlay(received_msg["data"]);
        }
    }
}

startSocket(document.currentScript.getAttribute('overlay'))

