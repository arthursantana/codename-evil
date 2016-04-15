var ResourceList = React.createClass({
   render: function () {
      if (this.props.planet.notEnoughWorkers)
         enoughWorkers = "notEnough";
      else
         enoughWorkers = "";

      if (this.props.planet.notEnoughCattle)
         enoughCattle = "notEnough";
      else
         enoughCattle = "";

      if (this.props.planet.notEnoughEnergy)
         enoughEnergy = "notEnough";
      else
         enoughEnergy = "";

      return (
         <div id="resourceList">
            <span className={enoughWorkers}>{this.props.planet.busyWorkers} / {this.props.planet.workers}</span><span className="resource workers" title="Workers"></span><br />
            <span className={enoughCattle}>{this.props.planet.busyCattle} / {this.props.planet.cattle}</span><span className="resource cattle" title="Cattle"></span><br />
            <span className={enoughEnergy}>{this.props.planet.busyEnergy} / {this.props.planet.energy}</span><span className="resource energy" title="Energy"></span><br />
            <span>{this.props.planet.obtanium}</span><span className="resource obtanium" title="Obtanium"></span>
         </div>
      );
   }
});
