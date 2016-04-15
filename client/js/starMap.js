var StarMap = React.createClass({
   drawPlanets: function () {
      if (!this.props.planets || !this.props.players)
         return;

      var ctx = this.refs.canvas.getContext('2d');

      ctx.font = "bold 10px Orbitron";
      ctx.textAlign = "center";
      for (var i = 0; i < this.props.planets.length; i++) {
         var p = this.props.planets[i];

         ctx.beginPath();
         ctx.arc(p.position[0], p.position[1], 5*p.r, 0, 2*Math.PI, false);
         if (p.ownerId == -1)
            ctx.fillStyle = '#ccc';
         else
            ctx.fillStyle = this.props.players[p.ownerId].color;
         ctx.lineWidth = 2;
         ctx.strokeStyle = '#fff';
         ctx.stroke();
         ctx.fill();
         ctx.fillStyle = "white";
         ctx.fillText(p.name,p.position[0],p.position[1]+5*p.r+20);
      }
   },

   drawShips: function () {
      if (!this.props.ships || !this.props.players)
         return;

      var ctx = this.refs.canvas.getContext('2d');

      ctx.font = "bold 10px Orbitron";
      ctx.textAlign = "center";
      for (var i = 0; i < this.props.ships.length; i++) {
         var s = this.props.ships[i];
         if (s.planetId != -1)
            continue;

         var r = 1;

         ctx.beginPath();
         ctx.arc(s.position[0], s.position[1], r, 0, 2*Math.PI, false);
         if (s.ownerId == -1)
            ctx.fillStyle = '#ccc';
         else
            ctx.fillStyle = this.props.players[s.ownerId].color;
         ctx.lineWidth = 1;
         ctx.strokeStyle = '#fff';
         ctx.stroke();
         ctx.fill();
         ctx.fillStyle = "white";
         ctx.fillText(s.name,s.position[0],s.position[1]+r+20);
      }
   },

   clickHandler: function (event) {
      x = event.pageX - this.refs.canvas.offsetLeft;
      y = event.pageY - this.refs.canvas.offsetTop;

      hasClickedAnyPlanet = false;
      for (var i = 0; i < this.props.planets.length; i++) {
         var p = this.props.planets[i];

         dx = x-p.position[0];
         dy = y-p.position[1];
         r = 5*p.r;

         if (dx*dx + dy*dy < r*r) {
            if (this.props.selectedShip == null) {
               hasClickedAnyPlanet = true;

               this.props.setSelectedPlanet(p);
            } else {
               var workers = Number(prompt("Crew size", "10000"));
               var cattle = Number(prompt("Cattle cargo size", "10000"));
               var obtanium = Number(prompt("Obtanium cargo size", "1000"));

               socket.send(JSON.stringify({
                  command: "setVoyage",
                  paramsSetVoyage: {
                     shipId: this.props.selectedShip.id,
                     destinationId: p.id,
                     workers: workers,
                     cattle: cattle,
                     obtanium: obtanium
                  }
               }));
            }
         }

         if (this.props.selectedShip != null)
            this.props.quitSetVoyageMode();
      }
   },

   render: function () {
      var voyageModeClass = "";

      if (this.props.selectedShip != null)
         voyageModeClass = "voyageMode"

      if (this.props.planets == null || this.props.players == null)
         return null;

      return <canvas ref="canvas" id="starMap" className={voyageModeClass} width="750" height="750" onClick={this.clickHandler}></canvas>;
   },

   componentDidUpdate: function () {
      var ctx = this.refs.canvas.getContext('2d');

      ctx.clearRect(0, 0, this.refs.canvas.width, this.refs.canvas.height);
      this.drawPlanets();
      this.drawShips();
   }
});
