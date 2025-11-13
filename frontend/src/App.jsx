import { useEffect, useState } from 'react';
import './App.css';
import { HashRouter, Route, Routes } from 'react-router';
import Settings from './pages/Settings';
import Events from './pages/Events';
import TADrawer from './components/Drawer';
import { createTheme, ThemeProvider } from '@mui/material';
import { useDispatch, useSelector } from 'react-redux';

import { GetConfig, GetInfo } from "../wailsjs/go/main/App";
import { EventsOn } from '../wailsjs/runtime/runtime';

import Overview from './pages/Overview';
import { setConfig } from './redux/configSlice';
import { setInfo } from './redux/infoSlice';
import Chat from './pages/Chat';
import { setEvents } from './redux/eventsSlice';
import { setChatMsgs } from './redux/chatMsgsSlice';

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    }
})

function App() {
    const dispatch = useDispatch();
    const cfg = useSelector(state => state.config);

    useEffect(() => {
        // load productinfo like version, name, author, etc.
        GetInfo().then(info => {
            dispatch(setInfo(info))
        });

        // config, loads once, then reloads on every change. source of truth is ALWAYS go/db, never JS
        GetConfig().then(cfg => dispatch(setConfig(cfg)));
        EventsOn('config_update', (cfg) => {
            dispatch(setConfig(cfg));
        })
    }, [])

    useEffect(() => {
        if(!cfg || !cfg.Settings) return;

        const maxEvents = cfg.Settings.eventListMaxItems ?? 100;
        const cancelEventListener = EventsOn('platform_event', (data) => {
            console.log("EVENT:", data);
            dispatch(setEvents({maxEvents, data}));
        })

        const maxMsgs = cfg.Settings.chatListMaxItems ?? 100; // TODO FIXME: implement this
        const cancelChatMessageListener = EventsOn('platform_chatMessage', (data) => {
            console.log("CHATMESSAGE:", data);
            dispatch(setChatMsgs({maxMsgs, data}))
        })

        return () => {
            cancelEventListener();
            cancelChatMessageListener();
        }
    }, [cfg])

    return (
        <ThemeProvider theme={darkTheme}>
            <TADrawer />
            <div id="App" style={{ marginLeft: 70, padding: 30 }}>
                <HashRouter>
                    <Routes>
                        <Route index element={<Overview />} />
                        <Route path="events" element={<Events />} />
                        <Route path="chat" element={<Chat />} />
                        <Route path="settings" element={<Settings />} />
                    </Routes>
                </HashRouter>
            </div>
        </ThemeProvider>
    )
}

export default App
