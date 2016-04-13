var StarMap = React.createClass({
   drawPlanets: function () {
      if (!this.props.planets || !this.props.players) {
         console.log("NO PLANETS OR NO PLAYERS!")
         return;
      }

      var ctx = this.refs.canvas.getContext('2d');

      ctx.clearRect(0, 0, this.refs.canvas.width, this.refs.canvas.height);

      if (this.props.planets == null) return;

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
            hasClickedAnyPlanet = true;
            this.props.setSelectedPlanet(p);
         }
      }
   },

   render: function () {
      if (this.props.planets == null || this.props.players == null)
         return null;

      return <canvas ref="canvas" id="starMap" width="750" height="750" onClick={this.clickHandler}></canvas>;
   },

   componentDidUpdate: function () {
      this.drawPlanets();
   }
});
