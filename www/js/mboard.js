// mboard.js

const REST_API          = "/api";
const WS_SUBSCRIBER     = "ws://10.0.1.20:8000/ws/subscribers";
const WS_GAME           = "ws://localhost:8000/ws/games";
const WS_MANAGER        = "ws://localhost:8000/ws/manager";

const API_PARAM_GAME_CONFIG = "gameConfig";

const PERIODS = ["1st", "2nd", "3rd", "4th"];

let ctl           = null;
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
    command("CLOCK_STOP");


  } else {

    stopButton();
    command("CLOCK_START");

  }

} // toggleClock


function clockStop() {

  playButton();
  command("CLOCK_STOP");

} // clockStop


function getPossession() {

  var away = document.getElementById("away");
  var home = document.getElementById("home");

  if(away === null || home === null) {
    return false;
  }

  if(away.getAttribute("class").includes("btn-info")) {
    return false;
  } else {
    return true;
  }

} // getPossession


function togglePossession() {

  var away = document.getElementById("away");
  var home = document.getElementById("home");

  // TODO: create a mechanism that can keep set of classes and append/pop
  if(getPossession()) {
    away.setAttribute("class", "btn btn-info rounded text-uppercase p-5 w-100 standard");
    home.setAttribute("class", "btn btn-outline-info rounded text-uppercase p-5 w-100 standard");

    command("POSSESSION_AWAY", null, {"stop": getClockState()});

  } else {

    away.setAttribute("class", "btn btn-outline-info rounded text-uppercase p-5 w-100 standard");
    home.setAttribute("class", "btn btn-info rounded text-uppercase p-5 w-100 standard");

    command("POSSESSION_HOME", null, {"stop": getClockState()});

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
    command("FINAL");
    window.location = "/home";
  }

} // endGame


function confirmEndPeriod() {

  clockStop();

  let ans = confirm("If you wish to the end period, changes can no longer be made.");

  if(ans) {
    command("PERIOD_UP");
  }

} // confirmEndPeriod


function gameFinal() {

  alert("Game has completed, additional actions cannot be performed.");

  window.location = "/home";

} // gameFinal


function command(cmd, step, meta) {

  ctl.send(JSON.stringify({
    "cmd": cmd,
    "step": step,
    "meta": meta
  }));

} // command


function listener(obj) {

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
      console.log(obj);
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
      playButton();
      togglePossession();
      break;

    case "GAME_FINAL":
      gameFinal();
      break;

    default:
      break;

  }

} // listener


function subscribe(id) {

  subscriber = new WebSocket(`${WS_SUBSCRIBER}/${id}`);

  subscriber.onopen = function(e) {
    console.log("Subscribed successfully.")
  }

  subscriber.onclose = function(e) {
    console.log("Subscription closed, game does not exist or has been completed.");
  }

  subscriber.onmessage = function(e) {

    var obj = JSON.parse(e.data);

    listener(obj);

  }

  subscriber.onerror = function(e) {
    console.log(e);
  }

} // subscribe


function gamectl(id) {

  ctl = new WebSocket(`${WS_GAME}/${id}`);

  ctl.onopen = function(e) {
    console.log("Game connection successful.");
    ctl.send(JSON.stringify({"cmd": "GAME_STATE"}));
  }

  ctl.onmessage = function(e) {

  }

  ctl.onerror = function(e) {
    console.log(e);
  }

} // gamectl
