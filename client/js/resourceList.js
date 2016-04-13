var ResourceList = React.createClass({
   render: function () {
      return (
         <div id="resourceList">
            {this.props.planet.population}<span className="resource pop"></span><br />
            {this.props.planet.cattle}<span className="resource vaca"></span><br />
            {this.props.planet.energy}<span className="resource joules"></span><br />
            {this.props.planet.obtanium}<span className="resource obtanium"></span>
         </div>
      );
   }
});
