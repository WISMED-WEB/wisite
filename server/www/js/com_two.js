import { getEmitter } from './mitt.js'

let emitter = getEmitter();

export default {

    setup() {
        let msgFromApp1 = Vue.ref("");

        // listen to all events
        // emitter.on('*', (type, e) => console.log(type, e));

        // listen to an event, 'app1' is sender name
        emitter.on('app1', e => {
            console.log('app2 received from app1:', e)
            msgFromApp1.value = e;
        });

        return {
            msgFromApp1,
        };
    },

    template: `      
        <h1>from app1: {{msgFromApp1}}</h1>
    `,
};