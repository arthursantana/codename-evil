var BuildingInterface = React.createClass({
   build: function(type) {
      socket.send(JSON.stringify({
         command: "build",
         paramsBuild: {
            type: type,
            planetId: this.props.planet.id,
            i: this.props.selectedI,
            j: this.props.selectedJ
         }
      }));

      this.props.unselectBuilding();
   },

   train: function(type) {
      socket.send(JSON.stringify({
         command: "train",
         paramsTrain: {
            type: type,
            planetId: this.props.planet.id,
            i: this.props.selectedI,
            j: this.props.selectedJ
         }
      }));
   },

   sellBuilding: function() {
      socket.send(JSON.stringify({
         command: "sellBuilding",
         paramsSellBuilding: {
            planetId: this.props.planet.id,
            i: this.props.selectedI,
            j: this.props.selectedJ
         }
      }));

      this.props.unselectBuilding();
   },

   render: function () {
      var b = null;

      if (this.props.selectedI != -1 && this.props.selectedJ != -1)
         b = this.props.planet.buildings[this.props.selectedI][this.props.selectedJ];

      if (b == null)
         return null;

      if (b.ticksUntilDone > 0) {
         var count = b.ticksUntilDone

         interfaceBody = []
         while (count > 0) {
            interfaceBody.push(<div key={count} className="stillBuilding"></div>);
            count--;
         }
      } else if (b.operational == false) {
         interfaceBody = null
      } else {
         if (b.type == "") {
            interfaceBody = (
               <div>
                  <div>
                     <div className="tile farm" title={buildingTitle("farm")} onClick={this.build.bind(this, "farm")} ></div>
                     Farm
                     <div className="buildingCosts">
                        1k <span className="resource workers" title="Workers"></span>
                        100<span className="resource energy" title="Energy"></span>
                        100 <span className="resource obtanium" title="Obtanium"></span>
                     </div>
                  </div>
                  <div>
                     <div className="tile generator" title={buildingTitle("generator")} onClick={this.build.bind(this, "generator")} ></div>
                     Generator
                     <div className="buildingCosts">
                        5k <span className="resource workers" title="Workers"></span>
                        0 <span className="resource energy" title="Energy"></span>
                        200 <span className="resource obtanium" title="Obtanium"></span>
                     </div>
                  </div>
                  <div>
                     <div className="tile lockheed" title={buildingTitle("lockheed")} onClick={this.build.bind(this, "lockheed")} ></div>
                     Lockheed
                     <div className="buildingCosts">
                        15k <span className="resource workers" title="Workers"></span>
                        500 <span className="resource energy" title="Energy"></span>
                        200 <span className="resource obtanium" title="Obtanium"></span>
                     </div>
                  </div>
                  <div>
                     <div className="tile vale" title={buildingTitle("vale")} onClick={this.build.bind(this, "vale")} ></div>
                     Vale
                     <div className="buildingCosts">
                        20k <span className="resource workers" title="Workers"></span>
                        2k <span className="resource energy" title="Energy"></span>
                        1k <span className="resource obtanium" title="Obtanium"></span>
                     </div>
                  </div>
                  <div>
                     <div className="tile nasa" title={buildingTitle("nasa")} onClick={this.build.bind(this, "nasa")} ></div>
                     Nasa
                     <div className="buildingCosts">
                        30k <span className="resource workers" title="Workers"></span>
                        1k <span className="resource energy" title="Energy"></span>
                        500 <span className="resource obtanium" title="Obtanium"></span>
                     </div>
                  </div>
               </div>
            );
         } else if (b.type == "nasa") {
            if (b.ticksUntilProductionDone == 0)
               buildingQueue = null;
            else {
               var ticksUntilProductionDone = [];
               for (var count = b.ticksUntilProductionDone; count > 0; count--) {
                  ticksUntilProductionDone.push(<div key={count} className="stillBuilding"></div>);
               }

               var queue = [];
               for (var i = 0; i < b.productionQueue.length; i++) {
                  newShipIcon = <div key={i} data-id={i} className={"ShipIcon " + b.productionQueue[i]} title={b.productionQueue[i]}></div>;
                  queue.push(newShipIcon);
               }

               buildingQueue = (
                  <div>
                     <br />
                     Queue: {queue}<br /><br />
                     {ticksUntilProductionDone}
                  </div>
               );
            }
            interfaceBody = (
               <div>
                  <div>
                     <div className="unitIcon colonizer" onClick={this.train.bind(this,"colonizer")}></div>
                     Colonizer
                     <div className="unitCosts">
                        50k <span className="resource workers" title="Workers"></span>
                        20k <span className="resource cattle" title="Cattle"></span>
                        5k <span className="resource obtanium" title="Obtanium"></span>
                     </div>
                  </div>
                  <div>
                     <div className="unitIcon trojan" onClick={this.train.bind(this,"trojan")}></div>
                     Trojan
                     <div className="unitCosts">
                        2k <span className="resource obtanium" title="Obtanium"></span>
                     </div>
                  </div>
                  {buildingQueue}
               </div>
            );
         } else if (b.type == "lockheed") {
            if (b.ticksUntilProductionDone == 0)
               buildingQueue = null;
            else {
               var ticksUntilProductionDone = [];
               for (var count = b.ticksUntilProductionDone; count > 0; count--) {
                  ticksUntilProductionDone.push(<div key={count} className="stillBuilding"></div>);
               }

               var queue = [];
               for (var i = 0; i < b.productionQueue.length; i++) {
                  newUnitIcon = <div key={i} data-id={i} className={"unitIcon " + b.productionQueue[i]} title={b.productionQueue[i]}></div>;
                  queue.push(newUnitIcon);
               }

               buildingQueue = (
                  <div>
                     <br />
                     Queue: {queue}<br /><br />
                     {ticksUntilProductionDone}
                  </div>
               );
            }
            interfaceBody = (
               <div>
                  <div>
                     <div className="unitIcon soldier" onClick={this.train.bind(this,"soldier")}></div>
                     Soldier Unit
                     <div className="unitCosts">
                        5k <span className="resource workers" title="Workers"></span>
                        500 <span className="resource obtanium" title="Obtanium"></span>
                     </div>
                     {buildingQueue}
                  </div>
               </div>
            );
         } else {
            interfaceBody = null;
         }
      }

      if (b.type != "" && b.type != "hq") {
         sellButton = <a onClick={this.sellBuilding}><i className="fa fa-trash-o" aria-hidden={true}></i></a>;
      } else {
         sellButton = null;
      }

      return (
         <div>
            <div className="unselector" onClick={this.props.unselectBuilding}></div>

            <div id="buildingInterface">
               <h2><span>{buildingTitle(b)}</span>{sellButton}</h2>

               {interfaceBody}
            </div>
         </div>
      );
   }
});
