import { useEffect, useState } from 'react';
import './App.css';
import { HashRouter, Route, Routes } from 'react-router';
import Settings from './pages/Settings';
import Overview from './pages/Overview';

function App() {

    useEffect(() => {
    }, [])

    return (
        <div id="App">
            <a href="#">Overview</a>
            <a href="#settings">Settings</a>
            <HashRouter>
                <Routes>
                    <Route index element={<Overview />} />
                    <Route path="settings" element={<Settings />} />
                </Routes>
            </HashRouter>
        </div>
    )
}

export default App
