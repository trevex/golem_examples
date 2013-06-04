/*

   Copyright 2013 Niklas Voss

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

(function(global) {

    if (global["WebSocket"]) {
        var seperator = " ",
            DefaultJSONProtocol = {
                unpack: function(data) {
                    var name = data.split(seperator)[0];
                    return [name, data.substring(name.length+1, data.length)];
                },
                unmarshal: function(data) {
                    return JSON.parse(data);
                },
                marshal: function(data) {
                    return JSON.stringify(data);
                },
                pack: function(name, data) {
                    return name + seperator + data;
                }
            };

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
            protocol: DefaultJSONProtocol,
            setProtocol: function(protocol) {
                this.protocol = protocol;
            },
            onClose: function(evt) {
                if (this.debug) {
                    console.log("golem: Connection closed!");
                }
                if (this.callbacks["close"]) this.callbacks["close"](evt);
            },
            onMessage: function(evt) {
                var [name, data] = this.protocol.unpack(evt.data);
                if (this.debug) {
                    console.log("golem: Received "+name+"-Event.");
                }
                if (this.callbacks[name]) {
                    var obj = this.protocol.unmarshal(data);
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
                this.ws.send(this.protocol.pack(name, this.protocol.marshal(data)));
            }

        }

        global.golem = {
            Connection: Connection
        };

    } else {

        console.warn("golem: WebSockets not supported!");

    }
})(this)