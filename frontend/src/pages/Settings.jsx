import { useEffect, useState } from 'react';
import { GetConfig, SetProviderSettings, SetSettings } from "../../wailsjs/go/main/App";
import { Accordion, AccordionSummary, AccordionDetails, Typography, FormGroup, Checkbox, FormControlLabel, TextField, FormControl, InputLabel, Select, MenuItem, Tooltip } from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Help from '@mui/icons-material/Help';
import { capitalizeFirstLetter } from '../helper';
import { useDebounce } from 'use-debounce';
import { useSelector } from 'react-redux';

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
        </div>
    )
}
function UISettings({ _settings }) {
    const [settings, setSettings] = useState(_settings);
    const [debouncedSettings] = useDebounce(settings, 500);

    useEffect(() => {
        SetSettings(debouncedSettings);
    }, [debouncedSettings])

    return <div style={{ textAlign: "left" }}>
        <FormControl sx={{width: 300}}>
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
        <Tooltip sx={{margin: 2}} title={<div>
            <p>How many items to show in the Events List.</p>
            <p>High amount of items can lead to lag & takes more resources.</p>
        </div>}>
            <Help />
        </Tooltip>


    </div>
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

export default Settings
