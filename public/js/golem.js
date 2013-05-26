/*

	golem - lightweight Go WebSocket-framework
    Copyright (C) 2013  Niklas Voss

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.

*/

(function(global) {

    if (global["WebSocket"]) {
        var seperator = " ";

        function Connection(addr, debug) {
            
            this.ws = new WebSocket("ws://"+addr);
            
            this.callbacks = {};
            
            this.debug = debug

            this.ws.onclose = this.onClose.bind(this);
            this.ws.onopen = this.onOpen.bind(this);
            this.ws.onmessage = this.onMessage.bind(this);
        }

        Connection.prototype = {
            constructor: Connection,
            onClose: function(evt) {
                if (this.debug) {
                    console.log("golem: Connection closed!");
                }
                if (this.callbacks["close"]) this.callbacks["close"](evt);
            },
            onMessage: function(evt) {
                var data = evt.data,
                    name = data.split(seperator, 1)[0];
                if (this.debug) {
                    console.log("golem: Received "+name+"-Event.");
                }
                if (this.callbacks[name]) {
                    var json = data.substring(name.length+1, data.length),
                        obj  = JSON.parse(json);
                    this.callbacks[name](obj);
                }
            },
            onOpen: function(evt) {
                if (this.debug) {
                    console.log("golem: Connection established!");
                }
                if (this.callbacks["open"]) this.callbacks["open"](evt);
            },
            on: function(name, callback) {
                this.callbacks[name] = callback;
            },
            emit: function(name, data) {
                this.ws.send(name+" "+JSON.stringify(data));
            }

        }

        global.golem = {
            Connection: Connection
        };

    } else {

        console.warn("golem: WebSockets not supported!");

    }
})(this)