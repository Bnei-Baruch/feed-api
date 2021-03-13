import React, { 
  useState, 
} from 'react';

import { 
  Button,
  Checkbox,
  Dropdown,
  List,
} from 'semantic-ui-react';

const SUGGESTERS = [
  "ContentTypeSuggester",
  "RoundRobinSuggester",
  "CompletionSuggester",
  "SortSuggester",
  "ContentUnitsSuggester",
  "CollectionSuggester",
  "LastClipsSameTagSuggester",
  "LastClipsSuggester",
  "LastCollectionSameSourceSuggester",
  "LastCongressSameTagSuggester",
  "LastContentTypesSameTagSuggester",
  "LastContentUnitsSameCollectionSuggester",
  "LastContentUnitsSuggester",
  "LastLessonsSameTagSuggester",
  "LastLessonsSuggester",
  "LastProgramsSameTagSuggester",
  "LastProgramsSuggester",
  "NextContentUnitsSameSourceSuggester",
  "PrevContentUnitsSameCollectionSuggester",
  "PrevContentUnitsSameSourceSuggester",
  "RandomContentTypesSuggester",
];

const HAS_ARGS = [
  "ContentTypeSuggester",
  "ContentUnitsSuggester",
  "CollectionSuggester",
  "LastCollectionSameSourceSuggester",
  "LastContentTypesSameTagSuggester",
  "NextContentUnitsSameSourceSuggester",
  "PrevContentUnitsSameSourceSuggester",
  "RandomContentTypesSuggester",
];

const HAS_SECOND_ARGS = [
  "RandomContentTypesSuggester",
];

const HAS_SPECS = [
  "ContentTypeSuggester",
  "RoundRobinSuggester",
  "CompletionSuggester",
  "SortSuggester",
];

// Collection Types
const CT_ARTICLES           = "ARTICLES";
const CT_BOOKS              = "BOOKS";
const CT_CHILDREN_LESSONS   = "CHILDREN_LESSONS";
const CT_CLIPS              = "CLIPS";
const CT_CONGRESS           = "CONGRESS";
const CT_DAILY_LESSON       = "DAILY_LESSON";
const CT_FRIENDS_GATHERINGS = "FRIENDS_GATHERINGS";
const CT_HOLIDAY            = "HOLIDAY";
const CT_LECTURE_SERIES     = "LECTURE_SERIES";
const CT_LESSONS_SERIES     = "LESSONS_SERIES";
const CT_MEALS              = "MEALS";
const CT_PICNIC             = "PICNIC";
const CT_SONGS              = "SONGS";
const CT_SPECIAL_LESSON     = "SPECIAL_LESSON";
const CT_UNITY_DAY          = "UNITY_DAY";
const CT_VIDEO_PROGRAM      = "VIDEO_PROGRAM";
const CT_VIRTUAL_LESSONS    = "VIRTUAL_LESSONS";
const CT_WOMEN_LESSONS      = "WOMEN_LESSONS";

// Content Unit Types
const CT_ARTICLE               = "ARTICLE";
const CT_BLOG_POST             = "BLOG_POST";
const CT_BOOK                  = "BOOK";
const CT_CHILDREN_LESSON       = "CHILDREN_LESSON";
const CT_CLIP                  = "CLIP";
const CT_EVENT_PART            = "EVENT_PART";
const CT_FRIENDS_GATHERING     = "FRIENDS_GATHERING";
const CT_FULL_LESSON           = "FULL_LESSON";
const CT_KITEI_MAKOR           = "KITEI_MAKOR";
const CT_LECTURE               = "LECTURE";
const CT_LELO_MIKUD            = "LELO_MIKUD";
const CT_LESSON_PART           = "LESSON_PART";
const CT_MEAL                  = "MEAL";
const CT_PUBLICATION           = "PUBLICATION";
const CT_RESEARCH_MATERIAL     = "RESEARCH_MATERIAL";
const CT_SONG                  = "SONG";
const CT_TRAINING              = "TRAINING";
const CT_UNKNOWN               = "UNKNOWN";
const CT_VIDEO_PROGRAM_CHAPTER = "VIDEO_PROGRAM_CHAPTER";
const CT_VIRTUAL_LESSON        = "VIRTUAL_LESSON";
const CT_WOMEN_LESSON          = "WOMEN_LESSON";

const COLLECTION_TYPES = [
  // Collection Types
  CT_ARTICLES,
  CT_BOOKS,
  CT_CHILDREN_LESSONS,
  CT_CLIPS,
  CT_CONGRESS,
  CT_DAILY_LESSON,
  CT_FRIENDS_GATHERINGS,
  CT_HOLIDAY,
  CT_LECTURE_SERIES,
  CT_LESSONS_SERIES,
  CT_MEALS,
  CT_PICNIC,
  CT_SONGS,
  CT_SPECIAL_LESSON,
  CT_UNITY_DAY,
  CT_VIDEO_PROGRAM,
  CT_VIRTUAL_LESSONS,
  CT_WOMEN_LESSONS,

];

const ALL_CONTENT_TYPES = [
  // Collection Types
  CT_ARTICLES,
  CT_BOOKS,
  CT_CHILDREN_LESSONS,
  CT_CLIPS,
  CT_CONGRESS,
  CT_DAILY_LESSON,
  CT_FRIENDS_GATHERINGS,
  CT_HOLIDAY,
  CT_LECTURE_SERIES,
  CT_LESSONS_SERIES,
  CT_MEALS,
  CT_PICNIC,
  CT_SONGS,
  CT_SPECIAL_LESSON,
  CT_UNITY_DAY,
  CT_VIDEO_PROGRAM,
  CT_VIRTUAL_LESSONS,
  CT_WOMEN_LESSONS,

  // Content Unit Types
  CT_ARTICLE,
  CT_BLOG_POST,
  CT_BOOK,
  CT_CHILDREN_LESSON,
  CT_CLIP,
  CT_EVENT_PART,
  CT_FRIENDS_GATHERING,
  CT_FULL_LESSON,
  CT_KITEI_MAKOR,
  CT_LECTURE,
  CT_LELO_MIKUD,
  CT_LESSON_PART,
  CT_MEAL,
  CT_PUBLICATION,
  CT_RESEARCH_MATERIAL,
  CT_SONG,
  CT_TRAINING,
  CT_UNKNOWN,
  CT_VIDEO_PROGRAM_CHAPTER,
  CT_VIRTUAL_LESSON,
  CT_WOMEN_LESSON,
].sort();

const splitKeyTail = (key) => {
  const parts = key.split('.');
  const tailParts = parts.slice(-1)[0].split('-');
  return [/*rest*/ parts.slice(0, -1).join('.'), /*tail index*/ Number(tailParts[0]), /*tail suggester*/ tailParts[1]];
};

const splitKeyHead = (key) => {
  const parts = key.split('.');
  const headParts = parts.slice(0)[0].split('-');
  return [/*head index*/ Number(headParts[0]), /*head suggester*/ headParts[1], /*rest*/ parts.slice(1).join('.')];
};

const SpecTree = (props) => {
  const {spec, onChange} = props;
  const [expanded, setExpanded] = useState([]);
  const [selected, setSelected] = useState('');

  const clickToggle = (key) => {
    const index = expanded.indexOf(key);
    if (index !== -1) {
      expanded.splice(index, 1);
    } else {
      expanded.push(key);
    }
    setExpanded(expanded.slice());
  }

  const SpecItem = (prefix, spec) => {
    const newPrefix = prefix ? `${prefix}-${spec.name}` : `0-${spec.name}`;
    const specExpanded = expanded.includes(newPrefix);
    return (
      <List key={newPrefix} selection>
        <List.Item active={selected === newPrefix}>
          {spec.specs && spec.specs.length ? <List.Icon name={specExpanded ? 'minus' : 'plus'} onClick={() => clickToggle(newPrefix)} /> : null}
          {(!spec.specs || !spec.specs.length) ? <div style={{'paddingRight': '1.8em', 'display': 'table-cell'}}></div> : null}
          <List.Content style={{'display': (!spec.specs || !spec.specs.length) ? 'table-cell' : undefined}} onClick={(e) => {e.stopPropagation(); setSelected(newPrefix);}}>
            <List.Header style={{'fontWeight': selected !== newPrefix ? 'normal' : undefined}}>
              {spec.name}{spec.args ? `[${spec.args.join(',')}]` : ''}
            </List.Header> 
            { specExpanded && spec.specs && spec.specs.map((child, index) => SpecItem(`${newPrefix}.${index}`, child)) }
          </List.Content>
        </List.Item>
      </List>
    );
  }

  const find = (selected, spec) => {
    if (selected.split('-')[1] === spec.name) {
      return spec;
    }
    const parts = selected.split('.');
    const index = Number(parts[1].split('-')[0]);
    const child = spec.specs[index];
    return find(parts.slice(1).join('.'), child);
  }

  const add = (suggester) => {
    if (!selected) {
      setSelected(`0-${suggester}`);
      onChange({name: suggester});
    } else {
      const selectedSpec = find(selected, spec);
      if (!selectedSpec.specs) {
        selectedSpec.specs = [];
      }
      selectedSpec.specs.push({name: suggester});
      onChange(spec);
    }
  }
  
  const remove = (currentSelected, currentSpec) => {
    if (currentSelected.split('-')[1] === currentSpec.name) {
      setSelected('');
      setExpanded([]);
      onChange(null);
      return;
    }
    const parts = currentSelected.split('.');
    const childIndex = parts[1].split('-')[0];
    if (parts.length === 2) {
      currentSpec.specs.splice(childIndex, 1);
      if (currentSpec.length === 0) {
        
      }
      onChange(spec);
      return;
    }
    return remove(parts.slice(1).join('.'), currentSpec.specs[childIndex]);
  }

  const fixRemoveSelectedExpanded = (removedKey) => {
    const [removeRest, removeTailIndex] = splitKeyTail(removedKey);
    const newExpanded = expanded.map((key) => {
      if (key.startsWith(removedKey)) {
        // remove.
        return null;
      }
      if (key.startsWith(removeRest)) {
        const keyTail = key.slice(removeRest.length);
        const [headIndex, headSuggester, rest] = splitKeyHead(keyTail);
        if (headIndex > removeTailIndex) {
          return [removeRest, `${headIndex+1}-${headSuggester}`, rest].join('.');
        }
      }
      return key;
    }).filter((key) => key);
    setExpanded(newExpanded);
    setSelected('');
  }

  const selectedSpec = (selected && find(selected, spec)) || null;
  const toggleArg = (contentType) => {
    if (!selectedSpec.args) {
      selectedSpec.args = [];
    }
    const index = selectedSpec.args.indexOf(contentType);
    if (index !== -1) {
      selectedSpec.args.splice(index, 1);
    } else {
      selectedSpec.args.push(contentType);
    }
    onChange(spec);
  };
  const addSecondArg = (tagUid) => {
    console.log(tagUid, selectedSpec);
    if (!selectedSpec.second_args) {
      selectedSpec.second_args = [];
    }
    const index = selectedSpec.second_args.indexOf(tagUid);
    if (index === -1) {
      selectedSpec.second_args.push(tagUid);
    }
    onChange(spec);
  }
  const setSecondArgs = (tagUids) => {
    console.log(tagUids, selectedSpec);
    if (!selectedSpec.second_args) {
      selectedSpec.second_args = [];
    }
    selectedSpec.second_args = tagUids;
    onChange(spec);
  }

  return (
    <div>
      <Button.Group>
        <Dropdown
          button
          options={SUGGESTERS.map(s => ({key: s, text: s, value: s}))}
          text='Add'
          onChange={(event, data) => add(data.value)}
          disabled={!!selected && !HAS_SPECS.includes(selected.split('.').slice(-1)[0].split('-')[1])}
          value={''}
        />
        <Button disabled={!selected} onClick={() => { remove(selected, spec); fixRemoveSelectedExpanded(selected); }}>Remove</Button>
        <Dropdown
          button
          scrolling
          closeOnChange={false}
          multiple
          text='Edit Content Types'
          disabled={!selected || !HAS_ARGS.includes(selected.split('.').slice(-1)[0].split('-')[1])}>
          <Dropdown.Menu>
            {ALL_CONTENT_TYPES.map(contentType => (
              <Dropdown.Item key={contentType + selectedSpec?.args?.includes(contentType)} onClick={(e) => e.stopPropagation()}>
                <Checkbox
                  checked={selectedSpec?.args?.includes(contentType)}
                  label={`${contentType} (${COLLECTION_TYPES.includes(contentType) ? 'collection' : 'content unit'})`}
                  onClick={() => toggleArg(contentType)}
                />
              </Dropdown.Item>
            ))}
          </Dropdown.Menu>
        </Dropdown>
        <Dropdown
          button
          options={selectedSpec?.second_args?.map((value) => ({text: value, value})) ?? []}
          placeholder='Add Tag Uids'
          search
          multiple
          selection
          allowAdditions
          value={selectedSpec?.second_args ?? []}
          onAddItem={(event, data) => addSecondArg(data.value)}
          onChange={(event, data) => setSecondArgs(data.value)}
          disabled={!selected || !HAS_SECOND_ARGS.includes(selected.split('.').slice(-1)[0].split('-')[1])}
        />
      </Button.Group>
      { spec && SpecItem('', spec, expanded) }
    </div>
  );
};


export default SpecTree;
