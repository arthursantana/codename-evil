var planetScale = 5;
var mouseX = 0;
var mouseY = 0;

var StarMap = React.createClass({
   drawPlanets: function () {
      if (!this.props.planets || !this.props.players)
         return;

      var ctx = this.refs.canvas.getContext('2d');

      ctx.font = "bold 10px Questrial";
      ctx.textAlign = "center";
      for (var i = 0; i < this.props.planets.length; i++) {
         var p = this.props.planets[i];

         r = planetScale*p.r;

         ctx.beginPath();
         ctx.arc(p.position[0], p.position[1], r, 0, 2*Math.PI, false);
         if (p.ownerId == -1)
            ctx.fillStyle = '#333';
         else
            ctx.fillStyle = this.props.players[p.ownerId].color;
         if (isInCircle(mouseX, mouseY, p.position[0], p.position[1], r)) {
            ctx.lineWidth = 3;
            ctx.strokeStyle = '#0f0';
            ctx.stroke();
         }
         ctx.fill();
         ctx.fillStyle = "white";
         ctx.fillText(p.name,p.position[0],p.position[1]+r+20);
      }
   },

   drawShips: function () {
      if (!this.props.ships || !this.props.players)
         return;

      var ctx = this.refs.canvas.getContext('2d');

      ctx.font = "bold 10px Questrial";
      ctx.textAlign = "center";
      for (var i = 0; i < this.props.ships.length; i++) {
         var s = this.props.ships[i];
         if (s.planetId != -1 || s.ownerId == -1)
            continue;

         var r = 1;

         ctx.beginPath();
         ctx.arc(s.position[0], s.position[1], r, 0, 2*Math.PI, false);
         ctx.fillStyle = this.props.players[s.ownerId].color;
         ctx.fill();
         ctx.fillStyle = "white";
         ctx.fillText(s.type,s.position[0],s.position[1]+r+20);
      }
   },

   mouseMoveHandler: function (event) {
      r = getMousePos(this.refs.canvas, event)
      mouseX = r.x;
      mouseY = r.y;
   },

   clickHandler: function (event) {
      x = event.pageX - this.refs.canvas.offsetLeft;
      y = event.pageY - this.refs.canvas.offsetTop;

      hasClickedAnyPlanet = false;
      for (var i = 0; i < this.props.planets.length; i++) {
         var p = this.props.planets[i];

         if (isInCircle(x, y, p.position[0], p.position[1], planetScale*p.r)) {
            if (this.props.selectedShip == null) {
               hasClickedAnyPlanet = true;

               this.props.setSelectedPlanet(p);
            } else {
               socket.send(JSON.stringify({
                  command: "setDestination",
                  paramsSetDestination: {
                     shipId: this.props.selectedShip.id,
                     destinationId: p.id
                  }
               }));
            }
         }

         if (this.props.selectedShip != null)
            this.props.quitSetDestinationMode();
      }
   },

   render: function () {
      var voyageModeClass = "";

      if (this.props.selectedShip != null) {
         voyageModeClass = "voyageMode"

         switch (this.props.selectedShip.type) {
            case "colonizer": voyageModeClass += " colonizer"; break;
            case "trojan": voyageModeClass += " trojan"; break;
         }
      }

      if (this.props.planets == null || this.props.players == null)
         return null;

      return <canvas ref="canvas" id="starMap" className={voyageModeClass} width="1300" height="750" onClick={this.clickHandler} onMouseMove={this.mouseMoveHandler}></canvas>;
   },

   componentDidUpdate: function () {
      var ctx = this.refs.canvas.getContext('2d');

      ctx.clearRect(0, 0, this.refs.canvas.width, this.refs.canvas.height);
      this.drawPlanets();
      this.drawShips();
   }
});
