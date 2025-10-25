import {useEffect, useState} from 'react';
import {GetConfig} from "../../wailsjs/go/main/App";

function Settings() {
    const [cfg, setCfg] = useState({});

    useEffect(() => {
        GetConfig().then(setCfg);
        console.log(cfg);
    }, [])

    return (
        <div id="App">
            <div dangerouslySetInnerHTML={{__html: cfg}}></div>
        </div>
    )
}

export default Settings
