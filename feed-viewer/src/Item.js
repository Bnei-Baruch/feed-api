import React from 'react';

import { 
  Button,
  Container,
  Icon,
  Segment,
} from 'semantic-ui-react';

import { CT_DAILY_LESSON, CT_VIDEO_PROGRAM } from './helpers/consts';
import { canonicalLink } from './helpers/links';

const Item = (props) => {
  const {id: mdbUid, content_type: contentType, name, content_units: contentUnits, film_date: filmDate} = props.item || {};
  const {suggester} = props.contentItem || {suggester: 'NotDefinedSuggester'};
  const toLink = canonicalLink(props.item || { id: mdbUid, content_type: contentType }, "he");
  console.log(props.item);
  console.log(props.contentItem);

	const renderCU = cu => {
		return (
		  <a
			key={cu.id}
			href={canonicalLink(cu, "he")}
			target="_blank"
            rel="noopener noreferrer"
		  >
			<Button basic size="tiny" className="link_to_cu">
			  {cu.name}
			</Button>
		  </a>
		);
	};

  if ([CT_DAILY_LESSON, CT_VIDEO_PROGRAM].includes(contentType)) {
    return (
      <Segment key={mdbUid}>
        <Container>
        <Container as="h3">
          <a className="search__link"
             href={toLink}
             target="_blank"
            rel="noopener noreferrer">
          { contentType === CT_DAILY_LESSON ?
            <span>שיעור בוקר {filmDate}</span> :
            <span>תוכנית {name}, פרק אחרון ב-{contentUnits[0].film_date}</span> }
          </a>
        </Container>
        {
        false &&
        <Container className="content">
          <span>{contentUnits.length}{' '}פרקים</span>
        </Container>
        }
        <div className="clear" />
        </Container>

        <Container className="content clear margin-top-8">
        {contentUnits.slice(0, 5).map(renderCU)}

        <a
          href={toLink}
          target="_blank"
                    rel="noopener noreferrer"
          style={{'whiteSpace': 'nowrap'}}
        >
          <Icon name="tasks" size="small" />
          &nbsp;
          {`הראה את כל ${contentUnits.length} החלקים`}
        </a>
        </Container>
        </Segment>
    );
  } else {
    return (
      <Segment key={mdbUid}>
        <Container>
        <h3>{filmDate}</h3>
        - {name} - {mdbUid} - {contentType} - {suggester}
        </Container>
      </Segment>
    );
  }
};


export default Item;
