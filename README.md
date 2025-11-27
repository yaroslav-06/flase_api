> [!CAUTION]
> This is a very early version. Default use will be over HTTP, which is not a secure connection. Thus, all messages can be read by anyone in between client and server. To have a secure connection the it should be updated from TCP to TLS.

# Flase API docs
This is an API, to schedule messages, and send them in different ways (but currently only telegram works).\
Designed for **Linux**, but should also work on Macos.

The inspiration behind the project was a TV-series, where the main character knew some infromation that endangered his life. So he programmed a server to release all that information in a specified time, if something happened to him.

I liked that idea, and wanted to recreate in a long time. And while I was thinking about I though of some additional features:

1) Multiple users can be connected to the same server. To use API directly, you should first login using _websocket send_ (hereafter referred to as an API call):
```json5
{
    r: "login",
    d: {
        username: "<YOUR USERNAME>",
        password: "<YOUR PASSWORD>",
    },
}

```
2) There can be infinitely many messages (hereafter called packets), and packets can be created using an API call (but only after login):
```json5
{
    r: "packet creator",
    d: {
        name: "<PACKAGE DISPLAY NAME>",
        deliveryTime: "<TIME OF DELIVERY>", //for example: (2025-11-21 05:50:50)
        pass: "<PACKAGE PASSCODE>",
        actions: [], //I will explain this later
        recievers: [
            {
                type: "telegram", // currently only telegram reciever works
                username: "<TELEGRAM USERNAME>", // for example @parolk06
                message: "<THE MESSAGE ITSELF>",
            }
        ],
    }
}
```
3) Each packet has a set of action that are available for execution (currently there are only _delete_, and _time changer_ actions):
```json5
{
    r: "packet creator",
    d: {
        //...
        actions: [
            {
                type: "delete",
                pass: "<ACTION PASSCODE>",
            },
            {
                type: "time changer",
                pass: "<ACTION PASSCODE>",
                duration: 1000, // by what duration (in minutes) to change the delivery time
            },
        ],
        //...
    }
}
```
4. Each packet can be shown on screen in the app with packet title, and execution time countdown, or through api with request (only after login):
```json5
{
    r: "get packet",
    d: {
        pass: "<PACKAGE PASSCODE>",
    }
}
```
5. After you display a packege you can execute any action when you enter that action's passcode, or with API call (only after login and get_package):
```json5
{
    r: "perform action",
    d: {
        pass: "<ACTION PASSCODE>"
    }
}
```

> [!NOTE]
>Motivation for existance of different actions was to prevent someone from forcibly obtaining a cancellation passcode.
 
> [!TIP]
> The particular example from the TV-show is very interesting, but not that common in the real life.\
> But I have found some other interesting use cases. For example, I use it as an after death message:\
> I create 2 packets. In first, I add all messages I want to send to everybody if something happens to me (information like passwords to my crypto accouns), and I set it to execute in the next year. In second packet I put remainder to myself to reschedule the first packet, and I set it to execute a week before the first packet.

## Installation

Install redis:

For example for debian based distributions:
```bash
sudo apt install redis
```

Then start redis server:
```bash
sudo systemctl start redis-server
sudo systemctl enable redis-server
systemctl status redis-server
```

Clone the repo:
```bash
git clone https://github.com/yaroslav-06/flase_api.git
cd flase_api
```

If you want telegram working you first have to configure [https://github.com/yaroslav-06/api_telegram_sender](api_telegram_sender).\
Then open internal/telegram/send.go, and set the const _your_api_key_, to be equal to _your_api_key_ from your [https://github.com/yaroslav-06/api_telegram_sender](api_telegram_sender):
```bash
vim internal/telegram/send.go
```
or
```bash
nano internal/telegram/send.go
```
Now this test should be successful:
```bash
go test ./internal/telegram
```

Now you can run the code (at first execution it will as for admin user and password):
```bash
go run cmd/main.go
```
