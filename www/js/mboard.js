// mboard.js

const REST_API          = "/api";
const WS_SUBSCRIBER     = "ws://localhost:8000/ws/subscribers";
const WS_SCORE          = "ws://localhost:8000/ws/scores";
const WS_CLOCK          = "ws://localhost:8000/ws/clocks";
const WS_MANAGER        = "ws://localhost:8000/ws/manager";

const API_PARAM_GAME_CONFIG = "gameConfig";

const PERIODS = ["1st", "2nd", "3rd", "4th"];

let ctl           = null;
let clockctl      = null;
let subscriber    = null;


function score(periods) {

  var keys = Object.keys(periods);

  var total = 0;

  for(var i = 0; i < keys.length; i++) {
    total = total + periods[keys[i]];
  }

  return total;

} // score


function toggleSlider(n, val) {

  document.getElementById(n).value = val;
  document.getElementById(`${n}Label`).innerText = val;

} // toggleSlider


function toggleSel(n, val) {

  var e = document.getElementById(n);

  e.defaultValue = val;

  for(var i = 0; i < e.children.length; i++ ) {

    if(e.children[i].children[0].value === val) {
      e.children[i].setAttribute("class", "rounded mr-2 btn btn-outline-info active");
    } else {
      e.children[i].setAttribute("class", "rounded mr-2 btn btn-outline-info");
    }

  }

} // toggleSel


function checkClockToggle() {

  var play = document.getElementById("play");
  var stop = document.getElementById("stop");

  if(play === null || stop === null) {
    return false;
  } else {
    return true;
  }

} // checkClockToggle


function playButton() {

  if(checkClockToggle()) {
    document.getElementById("play").setAttribute("class", "rounded p-5 w-100 btn btn-success");
    document.getElementById("stop").setAttribute("class", "d-none rounded p-5 w-100 btn btn-danger");
  }

} // playButton


function stopButton() {

  if(checkClockToggle()) {
    document.getElementById("play").setAttribute("class", "d-none rounded p-5 w-100 btn btn-success");
    document.getElementById("stop").setAttribute("class", "rounded p-5 w-100 btn btn-danger");
  }

} // stopButton


function gameClockToString(cur, mins) {

  var nCur = cur;

  if(cur.seconds > mins * 60) {
    nCur.seconds = cur.seconds % 60;
  }

  var delta   = mins * 60 - nCur.seconds;
  var ndelta  = delta - 1;
  var seconds = delta % 60;
  var minutes = Math.floor(delta/60);
  var tenths  = 10 - nCur.tenths;

  if(delta === 60) {

    if(minutes === 1) {

      if(tenths === 10) {
        return minutes + ":00";
      } else {
        return ndelta + "." + tenths;
      }

    } else {
      return minutes + ":59." + tenths;
    }

  } else if(minutes === 0) {

    if(ndelta === -1) {
      return "0.0";
    } else if(tenths === 10) {
      return delta + ".0";
    } else {
      return ndelta + "." + tenths;
    }

  } else if(seconds === 0) {
    return minutes + ":00";
  } else if(seconds < 10 && seconds >= 0) {
    return minutes + ":0" + seconds;
  } else {
    return minutes + ":" + seconds;
  }

} // gameClockToString


function shotClockToString(cur, shot) {

  if(cur.seconds > shot) {
    console.log("Shot clock current time exceeds shot clock boundaries.");
    return shot;
  } else {
    return shot - cur.seconds;
  }

} // shotClockToString


function getClockState() {

  if(checkClockToggle()) {

    var c = play.getAttribute("class");

    return !c.includes("d-none");

  } else {
    return false;
  }

} // getClockState


function toggleClock() {

  if(!getClockState()) {

    playButton();
    clockCommand("CLOCK_STOP");


  } else {

    stopButton();
    clockCommand("CLOCK_START");

  }

} // toggleClock


function clockStop() {

  playButton();
  clockCommand("CLOCK_STOP");

} // clockStop


function checkPossession() {

  var away = document.getElementById("awayPos");
  var home = document.getElementById("homePos");

  if(away === null || home === null) {
    return false;
  } else {
    return true;
  }

} // checkPossession


function updatePossession(team) {

  if(checkPossession()) {

    var away = document.getElementById("awayPos");
    var home = document.getElementById("homePos");

    if(team === "HOME") {
      away.setAttribute("class", "btn btn-outline-info rounded text-uppercase p-5 w-100 standard");
      home.setAttribute("class", "btn btn-info rounded text-uppercase p-5 w-100 standard");
    } else {
      away.setAttribute("class", "btn btn-info rounded text-uppercase p-5 w-100 standard");
      home.setAttribute("class", "btn btn-outline-info rounded text-uppercase p-5 w-100 standard");
    }

  }

} // updatePossession


function getPossession() {

  if(checkPossession()) {

    var away = document.getElementById("awayPos");

    if(away.getAttribute("class").includes("btn-info")) {
      return false;
    } else {
      return true;
    }

  } else {
    return false;
  }

} // getPossession


function setPossession(team) {

  var p = getPossession()

  if(p && team === "HOME") {
    return;
  }

  if(!p && team === "AWAY") {
    return;
  }

  if(team === "HOME") {
    clockCommand("POSSESSION_HOME", null, {"stop": getClockState()});
  } else {
    clockCommand("POSSESSION_AWAY", null, {"stop": getClockState()});
  }

} // setPossession


function togglePossession() {

  // TODO: create a mechanism that can keep set of classes and append/pop
  if(getPossession()) {

    //setPossession("AWAY");
    clockCommand("POSSESSION_AWAY", null, {"stop": getClockState()});

  } else {

    //setPossession("HOME");

    clockCommand("POSSESSION_HOME", null, {"stop": getClockState()});

  }

} // togglePossession


function callTimeout() {

  if(t === "HOME") {

    if(s === 1) {
      command("HOME_TIMEOUT");
    } else {
      command("HOME_TIMEOUT_CANCEL");
    }

  } else {

    if(s === 1) {
      command("AWAY_TIMEOUT");
    } else {
      command("AWAY_TIMEOUT_CANCEL");
    }

  }

} // callTimeout


function updateScore(team, val) {

  var home = document.getElementById("homeScore");
  var away = document.getElementById("awayScore");

  if(home === null || away === null) {
    return;
  }

  if(team === "HOME") {
    home.innerText = val;
  } else {
    away.innerText = val;
  }

} // updateScore


function updateFouls(team, val) {

  var home = document.getElementById("homeFouls");
  var away = document.getElementById("awayFouls");

  if(home === null || away === null) {
    return;
  }

  if(team === "HOME") {
    home.innerText = val;
  } else {
    away.innerText = val;
  }

} // updateFouls


function updateTeam(team, val) {

  var home = document.getElementById("home");
  var away = document.getElementById("away");

  if(home === null || away === null) {
    return;
  }

  if(team === "HOME") {
    home.innerText = val;
  } else {
    away.innerText = val;
  }

} // updateTeam


function updateTeamPos(team, val) {

  var homePos = document.getElementById("homePos");
  var awayPos = document.getElementById("awayPos");

  if(homePos === null || awayPos === null) {
    return;
  }

  if(team === "HOME") {
    homePos.innerText = val;
  } else {
    awayPos.innerText = val;
  }

} // updateTeamPos


function updateClock(gameCur, shotCur, minConf, shotConf) {

  var clock = document.getElementById("clock");
  var shot  = document.getElementById("shot");

  if(clock === null || shot === null) {
    return;
  }

  clock.innerText = gameClockToString(gameCur, minConf);

  // TODO: avoid if there's no shot clock
  shot.innerText = shotClockToString(shotCur, shotConf);

} // updateClock


function updatePeriod(v) {

  var period = document.getElementById("period");

  if(period === null) {
    return;
  }

  var p   = parseInt(v);
  var str = PERIODS[0];

  if(p > 3) {
    str = "OT" + (p - 3);
  } else {
    str = PERIODS[p];
  }

  period.innerText = str;

} // updatePeriod


function updateState(o) {

  if(o.final) {

    let ans = confirm("Game has ended, no further operations can be performed against this game.")

    if(ans) {
      window.location = "/home";
    }

  } else {

    updateScore("HOME", score(o.state.home.points));
    updateScore("AWAY", score(o.state.away.points));

    updateFouls("HOME", o.state.home.fouls);
    updateFouls("AWAY", o.state.away.fouls);

    updateTeam("HOME", o.state.home.name);
    updateTeam("AWAY", o.state.away.name);

    updateTeamPos("HOME", o.state.home.name);
    updateTeamPos("AWAY", o.state.away.name);

    updateClock(o.state.game, o.state.shot, o.state.settings.minutes,
      o.state.settings.shot);

    updatePeriod(o.state.period);

  }

} // updateState


function newGame() {

  var formData = new FormData();

  var j = JSON.stringify({
    "periods": parseInt(document.getElementById("periods").defaultValue),
    "minutes": parseInt(document.getElementById("minutes").value),
    "shot": parseInt(document.getElementById("shot").value),
    "timeouts": parseInt(document.getElementById("timeouts").value),
    "fouls": parseInt(document.getElementById("fouls").value),
    "home": document.getElementById("home").value,
    "away": document.getElementById("away").value
  });

  formData.append(API_PARAM_GAME_CONFIG, j);

  fetch(`${REST_API}/games`, {
    method: "post",
    body: formData
  })
  .then((response) => {
    if(response.ok) return response.text();
  })
  .then((data) => {
    console.log(data);
    window.location = `/clockctl/${data}`;
  })
  .catch((error) => {
    console.log(error);
  });

} // newGame


function endGame() {

  clockStop();

  let ans = confirm("Are you sure you want to end game?");

  if(ans) {
    clockCommand("FINAL");
    window.location = "/home";
  }

} // endGame


function checkPeriod() {

  var period  = document.getElementById("period");
  var away    = document.getElementById("awayScore");
  var home    = document.getElementById("homeScore");

  console.log(period.value);

  return true;

  //if(period.value === )
} // checkPeriod


function confirmEndPeriod() {

  clockStop();

  if(checkPeriod()) {

    let ans = confirm("If you wish to the end period, changes can no longer be made to that period.");

    if(ans) {
      clockCommand("PERIOD_UP");
    }

  }

} // confirmEndPeriod


function gameFinal() {

  alert("Game has completed, additional actions cannot be performed.");

  window.location = "/home";

} // gameFinal


function shotClockViolation() {

  playButton();
  togglePossession();

} // shotClockViolation


function clockCommand(cmd, step, meta) {

  clockctl.send(JSON.stringify({
    "cmd": cmd,
    "step": step,
    "meta": meta
  }));

} // clockCommand


function scoreCommand(cmd, step, meta) {

  ctl.send(JSON.stringify({
    "cmd": cmd,
    "step": step,
    "meta": meta
  }));

} // scoreCommand


function listener(obj) {

  console.log(obj);

  switch(obj.key) {

    case "HOME_SCORE":
      updateScore("HOME", obj.val);
      break;

    case "AWAY_SCORE":
      updateScore("AWAY", obj.val);
      break;

    case "HOME_FOUL":
      updateFouls("HOME", obj.val);
      break;

    case "AWAY_FOUL":
      updateFouls("AWAY", obj.val);
      break;

    case "GAME_STATE":
      updateState(obj);
      break;

    case "CLOCK":

      var j = JSON.parse(obj.val);

      updateClock(j.game, j.shot, j.minutes, j.shotclock);
      break;

    case "END_PERIOD":
      playButton();
      confirmEndPeriod();
      break;

    case "PERIOD":
      updatePeriod(obj.val);
      break;

    case "SHOT_VIOLATION":
      shotClockViolation();
      break;

    case "GAME_FINAL":
      gameFinal();
      break;

    case "POSSESSION_HOME":
      updatePossession("HOME");
      break;

    case "POSSESSION_AWAY":
      updatePossession("AWAY");
      break;

    default:
      break;

  }

} // listener


function getId() {

  var p = window.location.pathname.split("/");

  if(p.length > 0) {
    return p[p.length-1];
  } else {
    return "";
  }


} // getId


function subscribe() {

  var id = getId();

  subscriber = new WebSocket(`${WS_SUBSCRIBER}/${id}`);

  subscriber.onopen = function(e) {
    console.log("Subscribed successfully.");
  }

  subscriber.onclose = function(e) {
    console.log(e);
    //window.location = "/home";
  }

  subscriber.onmessage = function(e) {

    var obj = JSON.parse(e.data);

    listener(obj);

  }

  subscriber.onerror = function(e) {
    console.log("Subscribe error");
  }

} // subscribe


function scoreSocket() {

  var id = getId();

  ctl = new WebSocket(`${WS_SCORE}/${id}`);

  ctl.onopen = function(e) {
    console.log("Game connected successfully.");
    //ctl.send(JSON.stringify({"cmd": "GAME_STATE"}));
  }

  ctl.onmessage = function(e) {

  }

  ctl.onclose = function(e) {
    console.log(e);
    console.log("Unable to connect to game controller.");
  }

  ctl.onerror = function(e) {
    console.log(e);
  }

} // scoreSocket


function clockSocket() {

  var id = getId();

  clockctl = new WebSocket(`${WS_CLOCK}/${id}`);

  clockctl.onopen = function(e) {
    console.log("Clock connected successfully.");
    //clockctl.send(JSON.stringify({"cmd": "GAME_STATE"}));
  }

  clockctl.onmessage = function(e) {
    console.log(e.data);
  }

  clockctl.onclose = function(e) {
    alert("Unable to connect or connection lost, please try reconnecting.");
    //window.location = "/home";
  }

  clockctl.onerror = function(e) {
    console.log(e);
  }

} // clockSocket
