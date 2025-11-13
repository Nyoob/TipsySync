import { useEffect, useState } from 'react';
import { GetConfig, SetProviderSettings, SetSettings } from "../../wailsjs/go/main/App";
import { Accordion, AccordionSummary, AccordionDetails, Typography, FormGroup, Checkbox, FormControlLabel, TextField, FormControl, InputLabel, Select, MenuItem, Tooltip, Card } from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Help from '@mui/icons-material/Help';
import { capitalizeFirstLetter } from '../helper';
import { useDebounce } from 'use-debounce';
import { useSelector } from 'react-redux';
import { Info } from '@mui/icons-material';

function Settings() {
    const cfg = useSelector(state => state.config);
    if (!cfg.Providers || !cfg.Settings) { return null }

    return (
        <div>
            {/* <h1>General</h1> */}
            <h1>UI</h1>
            <UISettings _settings={cfg.Settings} />
            <h1>Platforms</h1>
            {Object.entries(cfg.Providers).map(([k, v], i) => {
                return <ProviderSettings provider={k} data={v} key={i} />
            })}
            <h1>Info</h1>
            <ProductInfo />
        </div>
    )
}

/*
 * UI SETTINGS
 */
function UISettings({ _settings }) {
    const [settings, setSettings] = useState(_settings);
    const [debouncedSettings] = useDebounce(settings, 500);

    useEffect(() => {
        SetSettings(debouncedSettings);
    }, [debouncedSettings])

    return <div style={{ textAlign: "left" }}>
        <FormControl sx={{ width: 300 }}>
            <InputLabel id="eventListMaxItemsLabel">Event-List Max Items</InputLabel>
            <Select
                labelId="eventListMaxItemsLabel"
                id="eventListMaxItems"
                value={settings.eventListMaxItems}
                label="Yarak"
                onChange={(e) => setSettings({ ...settings, eventListMaxItems: e.target.value })}
            >
                {["5", "10", "25", "50", "100", "250", "500", "1000"].map(x => {
                    return <MenuItem value={x}>{x}</MenuItem>
                })}
            </Select>
        </FormControl>
        <Tooltip sx={{ margin: 2 }} title={<div>
            <p>How many items to show in the Events List.</p>
            <p>High amount of items can lead to lag & takes more resources.</p>
        </div>}>
            <Help />
        </Tooltip>


    </div>
}

/* 
 * PROVIDER SETTINGS
 */
const apiTokenLabels = {
    chaturbate: "Events API Token URL",
    fansly: "Username",
    stripchat: "Username",
}
const shouldShowFetchInterval = {
    chaturbate: true,
    fansly: false,
    stripchat: false,
}
function ProviderSettings({ provider, data, key }) {
    const [settings, setSettings] = useState(data);
    const [debouncedSettings] = useDebounce(settings, 500);
    const [initialized, setInitialized] = useState(false) // to not call SetProviderSettings early

    useEffect(() => {
        if (!initialized) return;
        SetProviderSettings(provider, debouncedSettings);
    }, [debouncedSettings])

    useEffect(() => {
        setInitialized(true);
    })

    return <Accordion key={key}>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography component="span">{capitalizeFirstLetter(provider)}</Typography>
        </AccordionSummary>
        <AccordionDetails sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            <FormControlLabel label="Enabled?"
                control={<Checkbox checked={settings.Enabled} onChange={(e) => setSettings((s) => ({ ...s, Enabled: e.target.checked }))} />} />
            <TextField variant="standard" label={apiTokenLabels[provider]} value={settings.ApiToken}
                onChange={(e) => setSettings((s) => ({ ...s, ApiToken: e.target.value }))} />
            {shouldShowFetchInterval[provider] && <TextField variant="standard" label="Fetch Interval"
                slots={{ type: 'number' }} slotsProps={{ input: { min: 1, max: 120, step: 1 } }}
                value={settings.FetchInterval}
                onChange={(e) => setSettings((s) => ({ ...s, FetchInterval: parseInt(e.target.value) }))} />}
        </AccordionDetails>
    </Accordion>
}

/*
 * TOY CONTROL - NOT IMPLEMENTED YET
 */
function ToyControl() {
    return <div>
        {/* Here, expandable sections like providers */}
        {/* Inside those, add/removeable levels by native tips (converted by us to TipValueInDollars), also ranges */}
        {/* Setting these changes settings in DB */}
        {/* In Go, create new ToyProvider Service that waits for tipevents (from eventHandler) */}
        {/* These should then send signals to toys depending on configured levels */}

        {/* Likely a good idea to create bots for each platform aswell, idk tho yet */}
        {/* To send chat msg's like "toy connected" or "reacting to xyz" */}
        {/* good alternative would be stream overlay, idk if the default lovense overlay for example works */}
        {/* but doing our own overlay probably makes sense for other toys aswell. */}
    </div>
}

/*
 * PRODUCT INF
 */
function ProductInfo() {
    const info = useSelector(state => state.info);
    console.log("INFO:", info);

    if (!info || !info.name) return;

    return <Card>
        <div style={{ display: 'flex', padding: "32px 0"}}>
            <div style={{flex: 1, borderRight: "1px solid white"}}>
                <h3>App</h3>
                <p><b>App:</b> {info.name}</p>
                <p><b>Version:</b> {info.info.productVersion} ({info.info.build})</p>
                <p><b>Github:</b> {info.info.github}</p>
            </div>
            <div style={{flex: 1, borderLeft: "1px solid white"}}>
                <h3>Author</h3>
                <p><b>Made by:</b> {info.author.name}</p>
                <p><b>My Github:</b> {info.author.github}</p>
            </div>
        </div>

    </Card>
}

export default Settings
