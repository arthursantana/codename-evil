var BuildingInterface = React.createClass({
   build: function(type) {
      var i = this.props.selectedI;
      var j = this.props.selectedJ;

      socket.send(JSON.stringify({
         command: "build",
         paramsBuild: {
            type: type,
            planetId: this.props.planet.id,
            i: i,
            j: j
         }
      }));

      this.props.unselectBuilding();
   },

   render: function () {
      var b = this.props.selectedBuilding;

      if (b == null)
         return null

      if (b.type == "") {
         interfaceBody = (
            <div>
               <div>
                  <div className="tile farm" title={buildingTitle("farm")} onClick={this.build.bind(this, "farm")} ></div>
                  Farm
               </div>
               <div>
                  <div className="tile generator" title={buildingTitle("generator")} onClick={this.build.bind(this, "generator")} ></div>
                  Generator
               </div>
               <div>
                  <div className="tile vale" title={buildingTitle("vale")} onClick={this.build.bind(this, "vale")} ></div>
                  Vale
               </div>
               <div>
                  <div className="tile nasa" title={buildingTitle("nasa")} onClick={this.build.bind(this, "nasa")} ></div>
                  Nasa
               </div>
            </div>
         );
      } else {
         interfaceBody = null;
      }

      return (
         <div>
            <div className="unselector" onClick={this.props.unselectBuilding}></div>

            <div id="buildingInterface">
               <h2>{buildingTitle(b)}</h2>

               {interfaceBody}
            </div>
         </div>
      );
   }
});
