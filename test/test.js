var delay = ms => new Promise(res => setTimeout(res, ms))

var ws = new WebSocket("ws://localhost:14539")
await delay(300)
ws.onmessage = (e) => console.log(e.data)

var log = {
    r: "login",
    d: {
        username: "parolk",
        password: "qwerty",
    }
}

console.log(log)
ws.send(JSON.stringify(log))
await delay(300)

var cre = {
    r: "packet creator",
    d: {
        name: "test send packet",
        deliveryTime: "2025-11-21 05:50:50",
        pass: "tg10",
        actions: [
            {
                type: "delete",
                pass: "okok",
            },
            {
                type: "time changer",
                pass: "termint",
                duration: 1000
            },
        ],
        recievers: [
            {
                type: "telegram",
                username: "@parolk06",
                message: "on timout message backup",
            }
        ],
    }
}

console.log(cre)
ws.send(JSON.stringify(cre))
await delay(300)

var lod = {
    r: "get packet",
    d: {
        pass: "monto2"
    }
}

// console.log(lod)
// ws.send(JSON.stringify(lod))
// await delay(300)

// var act = {
//     r: "perform action",
//     d: {
//         pass: "termint"
//     }
// }
//
// console.log(act)
// ws.send(JSON.stringify(act))
// await delay(300)

var act = {
    r: "perform action",
    d: {
        pass: "okok"
    }
}

// console.log(act)
// ws.send(JSON.stringify(act))
// await delay(300)
//
// console.log(lod)
// ws.send(JSON.stringify(lod))
// await delay(300)
