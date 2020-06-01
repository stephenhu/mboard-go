var cs = null;
var ss = null;


function newGame() {

  var data = new FormData();

  var j = JSON.stringify({
    periods: parseInt(document.getElementById("periods").value),
    minutes: parseInt(document.getElementById("minutes").value),
    shot: parseInt(document.getElementById("shot").value),
    timeouts: parseInt(document.getElementById("timeouts").value),
    fouls: parseInt(document.getElementById("fouls").value),
    home: document.getElementById("home").value,
    away: document.getElementById("away").value
  });

  console.log(j);

  data.append("gameConfig", j);

  fetch("/api/games", {
    method: "post",
    body: data
  })
  .then((response) => {
    if(response.ok) return response.text();
  })
  .then((data) => {
    console.log(data);
    var gid = document.getElementById("gameId");
    gid.innerText = data;
    gid.value = data;
  })
  .catch((error) => {
    console.log(error);
  });


} // newGame


function subscribeListener(obj) {

  switch(obj.key) {
    case "HOME_SCORE":
      document.getElementById("homeScore").innerText = obj.val;
      break;

    case "AWAY_SCORE":
      document.getElementById("awayScore").innerText = obj.val;
      break;

    case "HOME_TIMEOUT":
      //document.getElementById("homeScore").value = obj.val
      break;

    case "AWAY_TIMEOUT":
      //updateTimeouts(AWAY, obj.val);
      break;

    case "HOME_FOUL":
      //updateFouls(HOME, obj.val);
      break;

    case "AWAY_FOUL":
      //updateFouls(AWAY, obj.val);
      break;

    case "CLOCK":

      var j = JSON.parse(obj.val);

      //updateClock(j.game, j.shot, j.minutes, j.shotclock);
      document.getElementById("gameClock").innerText = obj.val;
      break;

    case "POSSESSION_HOME":
      //updatePossession(HOME);
      break;

    case "POSSESSION_AWAY":
      //updatePossession(AWAY);
      break;

    case "PERIOD":
      //updatePeriod(obj.val);
      break;

    case "GAME_STATE":
      //updateDisplay(obj.state);
      break;

    default:
      break;

  }

} // subscribeListener


function control(id) {

  cs = new WebSocket("ws://127.0.0.1:8000/ws/games/" + id);

  cs.onmessage = function(e) {

    var o = JSON.parse(e.data);

    console.log(o);

  }

  cs.onerror = function(e) {
    console.log(e);
  }

  cs.onopen = function(e) {
    console.log("connected to game " + id);
  }

  cs.onclose = function(e) {
    console.log("closed connection " + id);
  }

} // control


function subscribe(id) {

  ss = new WebSocket("ws://127.0.0.1:8000/ws/subscribers/" + id);

  ss.onmessage = function(e) {

    var obj = JSON.parse(e.data);

    subscribeListener(obj);

  }

  ss.onerror = function(e) {
    console.log(e);
  }

  ss.onopen = function(e) {
    console.log("connected to game " + id + " as subscriber");
  }

  ss.onclose = function(e) {
    console.log("closed subscriber connection " + id);
  }

} // subscribe


function connectScoreboard() {

  var id = document.getElementById("gameId").value;

  subscribe(id);
  control(id);

} // connectScoreboard


function sendCommand(c, s) {

  console.log(c);
  console.log(s);

  if(s === "") {
    cs.send(JSON.stringify({"cmd": c}));
  } else {
    cs.send(JSON.stringify({"cmd": c, "step": s}));
  }

} // sendCommand
