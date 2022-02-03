import React, { 
  useEffect,
  useState, 
} from 'react';

import { 
  Button,
  Checkbox,
  Grid,
  Input,
  Segment,
} from 'semantic-ui-react';

import { scan } from './api'

const COLUMNS = [
  'id',
  'created_at',
  'user_id',
  'user_agent',
  'namespace',
  'client_event_id',
  'client_event_type',
  'client_flow_id',
  'client_flow_type',
  'client_session_id',
  'data',
];

const JSON_COLUMN = [
  'data',
];

const DEFAULT_COLUMNS_VIEW = {
	client_event_id: false,
	client_event_type: true,
	client_flow_id: false,
	client_flow_type: true,
	client_session_id: false,
	created_at: true,
	data: true,
	id: false,
	namespace: false,
	user_agent: false,
	user_id: false,
};

const DEFAULT_LIMIT = 1000;

const Chronicles = (props) => {
  const [id, setId] = useState(new URLSearchParams(window.location.search).get('id'));
  const [limit, setLimit] = useState(Number(new URLSearchParams(window.location.search).get('limit') || DEFAULT_LIMIT));
  const [eventTypes, setEventTypes] = useState((new URLSearchParams(window.location.search).get('event_types') || '').split(',').filter(a => a));
  const [userIds, setUserIds] = useState((new URLSearchParams(window.location.search).get('user_ids') || '').split(',').filter(a => a));
  const [namespaces, setNamespaces] = useState((new URLSearchParams(window.location.search).get('namespaces') || '').split(',').filter(a => a));
  const [keycloak, setKeycloak] = useState(!!(new URLSearchParams(window.location.search).get('keycloak')));


  const [isOpen, setIsOpen] = useState(true);
  // COnst [columnView, setColumnView] = useState(COLUMNS.reduce((view, column) => Object.assign(view, {[column]: true}), {}));
  const [columnView, setColumnView] = useState(DEFAULT_COLUMNS_VIEW);

  const [entries, setEntries] = useState([]);
  const [scanError, setScanError] = useState('');

  const scanClicked = (clear) => {
    setScanError('');
    scan(id, limit, eventTypes, userIds, namespaces, keycloak).then((newEntries) => {
      console.log(entries);
      setScanError('');
			setEntries(clear ? newEntries : entries.concat(newEntries));
			const newId = newEntries.slice(-1)[0].id || '';
			console.log('setting id', newId);
			setId(newId);
    }).catch((error) => {
      setScanError(String(error));
    });
  };
  useEffect(scanClicked, []);

  const checkColumnView = (name) => setColumnView({...columnView, [name]: !columnView[name]})
  // console.log(columnView);
  return (
    <Grid columns={2}>
        <Grid.Column width={isOpen ? 5 : 1}>
          <Segment style={{'direction': 'ltr', 'minHeight': '60px'}}>
            <Button onClick={() => setIsOpen(!isOpen)}
                    icon={isOpen ? 'right arrow' : 'left arrow'} circular compact floated='right'
                    style={{'marginRight': '-6px'}}/>
            <div style={{'display': isOpen ? 'block' : 'none'}}>,
              <h2 style={{'display': 'inline'}}>Scan</h2>
              <Segment textAlign='left'>
                <table>
                  <tbody>
                    <tr><td>Id:</td><td><Input placeholder='Id' value={id} onChange={(event, data) => setId(data.value)} /></td></tr>
										<tr><td>Limit:</td><td><Input placeholder='Limit' value={limit} onChange={(event, data) => setLimit(Number(data.value))} /></td></tr>
										<tr><td>EventTypes:</td><td><Input placeholder='EventTypes' value={eventTypes} onChange={(event, data) => setEventTypes(data.value.split(','))} /></td></tr>
										<tr><td>UserIds:</td><td><Input placeholder='UserIds' value={userIds} onChange={(event, data) => setUserIds(data.value.split(','))} /></td></tr>
										<tr><td>Namespaces:</td><td><Input placeholder='Namespaces' value={namespaces} onChange={(event, data) => setNamespaces(data.value.split(','))} /></td></tr>
										<tr><td>Keycloak:</td><td><Checkbox checked={keycloak} onChange={() => setKeycloak(!keycloak)} /></td></tr>
                    <tr><td>View</td><td>
											<table>
												<tbody>
												{COLUMNS.map(column =>
													<tr key={column}><td>
														<Checkbox key={column} checked={columnView[column]} onChange={() => checkColumnView(column)} label={column} />
													</td></tr>
												)}
												</tbody>
											</table>
                    </td></tr>
                  </tbody>
                </table>
              </Segment>
            </div>
          </Segment>
        </Grid.Column>
        <Grid.Column width={isOpen ? 11 : 15}>
          <Segment style={{'direction': 'ltr', 'textAlign': 'left', overflow: 'auto', maxHeight: '80vh'}}>
            <table>
              <thead>
                <tr>
                  {COLUMNS.map(column => columnView[column] && <th>{column}</th>)}
                </tr>
              </thead>
              <tbody>
                {entries.map((item, index) => {
                  return (
                    <tr key={item.id}>
                      {COLUMNS.map(column => {
                        if (columnView[column]) {
                          return (
                            <td style={{'whiteSpace': 'nowrap'}}>
															{JSON_COLUMN.includes(column) ? JSON.stringify(item[column]) : item[column]}
														</td>
                          );
                        }
                        return null;
                      })}
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </Segment>
          <Segment  style={{'direction': 'ltr'}}>
            <Button onClick={() => scanClicked(true)}>Scan</Button>
            <Button onClick={() => scanClicked(false)}>More</Button>
						Total: {entries.length}
            <br />
            {scanError && <span style={{color: 'red'}}>{scanError}</span>}
          </Segment>
        </Grid.Column>
    </Grid>
  );
}

export default Chronicles;
