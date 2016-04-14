var ResourceList = React.createClass({
   render: function () {
      return (
         <div id="resourceList">
            {this.props.planet.population}<span className="resource pop" title="Population"></span><br />
            {this.props.planet.cattle}<span className="resource vaca" title="Cattle"></span><br />
            {this.props.planet.energy}<span className="resource joules" title="Energy"></span><br />
            {this.props.planet.obtanium}<span className="resource obtanium" title="Obtanium"></span>
         </div>
      );
   }
});
