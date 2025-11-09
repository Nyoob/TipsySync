# Websocket
WIP...

Connect using `ws://localhost:6769/ws`.

Check [events.go](/internal/events/events.go) to see available events.
Events are sent as a `WrappedEvent{Platform string, Type string, Event Event}`
All websocket msg's are JSON encoded. An empty `{}` keepalive is sent every few seconds.
Platform can be any of the currently supported platforms, all lowercase.
Type can be one of: `[tip, follow, unfollow, subscribe, chatMessage]`.
Event data changes depending on Type.
