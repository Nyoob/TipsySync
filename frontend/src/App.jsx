import { useEffect, useState } from 'react';
import './App.css';
import { HashRouter, Route, Routes } from 'react-router';
import Settings from './pages/Settings';
import Events from './pages/Events';
import TADrawer from './components/Drawer';
import { createTheme, ThemeProvider } from '@mui/material';
import { useDispatch, useSelector } from 'react-redux';

import { GetConfig } from "../wailsjs/go/main/App";
import { EventsOn } from '../wailsjs/runtime/runtime';

import Overview from './pages/Overview';
import { setConfig } from './redux/configSlice';

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    }
})

function App() {
    const dispatch = useDispatch();
    const cfg = useSelector(state => state.config);
    const [events, setEvents] = useState([]);

    useEffect(() => {
        GetConfig().then(cfg => dispatch(setConfig(cfg)));

        EventsOn('config_update', (cfg) => {
            dispatch(setConfig(cfg));
        })
    }, [])

    useEffect(() => {
        if(!cfg || !cfg.Settings) return;

        const maxEvents = cfg.Settings.eventListMaxItems
        const cancelEventListener = EventsOn('platform_event', (data) => {
            console.log(data);
            setEvents(e => {
                var newE = [data, ...e];
                if (newE.length > maxEvents) {
                    newE.length = maxEvents;
                }
                return newE;
            });
        })
        const cancelChatMessageListener = EventsOn('platform_chatMessage', (data) => {
            console.log("CHATMESSAGE:", data);
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
                        <Route path="events" element={<Events events={events} />} />
                        <Route path="settings" element={<Settings />} />
                    </Routes>
                </HashRouter>
            </div>
        </ThemeProvider>
    )
}

export default App
