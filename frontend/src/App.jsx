import { useEffect, useState } from 'react';
import './App.css';
import { HashRouter, Route, Routes } from 'react-router';
import Settings from './pages/Settings';
import Overview from './pages/Overview';
import TADrawer from './components/Drawer';
import { createTheme, ThemeProvider } from '@mui/material';

import { EventsOn } from '../wailsjs/runtime/runtime';

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    }
})

function App() {
    const [events, setEvents] = useState([]);

    useEffect(() => {
        EventsOn('platform_event', (data) => {
            setEvents(e => {
                var newE = [data, ...e];
                if (newE.length > 100) {
                    newE.length = 100;
                }
                return newE;
            });
        })
    }, [])

    return (
        <ThemeProvider theme={darkTheme}>
            <TADrawer />
            <div id="App" style={{ marginLeft: 70, padding: 30 }}>
                <HashRouter>
                    <Routes>
                        <Route index element={<Overview events={events} />} />
                        <Route path="settings" element={<Settings />} />
                    </Routes>
                </HashRouter>
            </div>
        </ThemeProvider>
    )
}

export default App
