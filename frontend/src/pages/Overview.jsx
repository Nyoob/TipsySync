import { Card } from '@mui/material';
import { useEffect, useState } from 'react';
import { Responsive, WidthProvider } from "react-grid-layout";
import "react-grid-layout/css/styles.css";
import "react-resizable/css/styles.css";
import Events from './Events';
import Chat from './Chat';

const ResponsiveGridLayout = WidthProvider(Responsive);

export default function Overview({ }) {

  return <div>
    <ResponsiveGridLayout
      rowHeight={30}
      breakpoints={{ md: 800, xs: 0 }}
      cols={{ md: 12, xs: 6 }}
      margin={[32, 32]}
      width="100%">
      <div key="chat" data-grid={{ x: 6, y: 0, w: 6, h: 12, minW: 2 }}>
        <Card style={styles.card} className="scrollable">
          <Chat />
        </Card>
      </div>
      <div key="events" data-grid={{ x: 0, y: 0, w: 6, h: 12, minW: 2 }}>
        <Card style={styles.card} className="scrollable">
          <Events />
        </Card>
      </div>
    </ResponsiveGridLayout>
  </div>
}

const styles = {
  card: {
    padding: "20px",
    height: "calc(100% - 40px)",
  }
}
