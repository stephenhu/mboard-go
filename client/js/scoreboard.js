class Scoreboard {

  constructor(periods, minutes, shot, timeouts, fouls, away, home) {

    this.periods    = periods;
    this.minutes    = minutes;
    this.shot       = shot;
    this.timeouts   = timeouts;
    this.fouls      = fouls;
    this.away       = away;
    this.home       = home;

    this.id         = this.newGame();

  } // constructor


  newGame() {

    fetch("/api/games", {
      method: "post",
      body: data
    })
    .then((response) => {
      if(response.ok) return response.text();
    })
    .then((data) => {
      console.log(data);
      //var gid = document.getElementById("gameId");
      //gid.innerText = data;
      //gid.value = data;
      return data;
    })
    .catch((error) => {
      console.log(error);
    });

  }

} // Scoreboard

var s = new Scoreboard(4, 12, 24, 3, 10, "red", "blue");
