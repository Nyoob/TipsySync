import { useEffect, useState } from 'react';
import { GetConfig, SetProviderSettings } from "../../wailsjs/go/main/App";
import { Accordion, AccordionSummary, AccordionDetails, Typography, FormGroup, Checkbox, FormControlLabel, TextField } from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { capitalizeFirstLetter } from '../helper';
import { useDebounce } from 'use-debounce';

function Settings() {
    const [cfg, setCfg] = useState({});

    useEffect(() => {
        GetConfig().then(setCfg);
        console.log(cfg);
    }, [])

    if (!cfg.Providers) { return null }

    return (
        <div>
            {/* <h1>General</h1> */}
            <h1>Platforms</h1>
            {Object.entries(cfg.Providers).map(([k, v], i) => {
                return <ProviderSettings provider={k} data={v} key={i} />
            })}
        </div>
    )
}

function ProviderSettings({ provider, data, key }) {
    const [settings, setSettings] = useState(data);
    const [debouncedSettings] = useDebounce(settings, 500);

    useEffect(() => {
        SetProviderSettings(provider, debouncedSettings);
    }, [debouncedSettings])

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
    fansly: "Room ID"
}
const shouldShowFetchInterval = {
    chaturbate: true,
    fansly: false
}

export default Settings
