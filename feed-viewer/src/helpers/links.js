import { canonicalCollection } from './utils';

import {
  BLOG_ID_LAITMAN_CO_IL,
  BLOG_ID_LAITMAN_COM,
  BLOG_ID_LAITMAN_ES,
  BLOG_ID_LAITMAN_RU,
  COLLECTION_EVENTS_TYPE,
  COLLECTION_LESSONS_TYPE,
  COLLECTION_PROGRAMS_TYPE,
  COLLECTION_PUBLICATIONS_TYPE,
  CT_ARTICLE,
  CT_ARTICLES,
  CT_BLOG_POST,
  CT_CLIP,
  CT_CLIPS,
  CT_CONGRESS,
  CT_DAILY_LESSON,
  CT_EVENT_PART,
  CT_FRIENDS_GATHERING,
  CT_FRIENDS_GATHERINGS,
  CT_FULL_LESSON,
  CT_HOLIDAY,
  CT_LECTURE,
  CT_LECTURE_SERIES,
  CT_LESSON_PART,
  CT_LESSONS_SERIES,
  CT_MEAL,
  CT_MEALS,
  CT_PICNIC,
  CT_SPECIAL_LESSON,
  CT_UNITY_DAY,
  CT_VIDEO_PROGRAM,
  CT_VIDEO_PROGRAM_CHAPTER,
  CT_VIRTUAL_LESSON,
  CT_VIRTUAL_LESSONS,
  CT_WOMEN_LESSON,
  CT_WOMEN_LESSONS,
  CT_KTAIM_NIVCHARIM,
  EVENT_TYPES,
  UNIT_EVENTS_TYPE,
  UNIT_LESSONS_TYPE,
  UNIT_PROGRAMS_TYPE,
  UNIT_PUBLICATIONS_TYPE,
} from './consts';

const blogNames = new Map([
  [BLOG_ID_LAITMAN_RU, 'laitman-ru'],
  [BLOG_ID_LAITMAN_COM, 'laitman-com'],
  [BLOG_ID_LAITMAN_ES, 'laitman-es'],
  [BLOG_ID_LAITMAN_CO_IL, 'laitman-co-il'],
]);

const mediaPrefix = new Map([
  [CT_LESSON_PART, '/lessons/cu/'],
  [CT_LECTURE, '/lessons/cu/'],
  [CT_VIRTUAL_LESSON, '/lessons/cu/'],
  [CT_WOMEN_LESSON, '/lessons/cu/'],
  [CT_BLOG_POST, '/lessons/cu/'],
  // [CT_CHILDREN_LESSON, '/lessons/cu/'],
  [CT_KTAIM_NIVCHARIM, '/lessons/cu/'],
  [CT_VIDEO_PROGRAM_CHAPTER, '/programs/cu/'],
  [CT_CLIP, '/programs/cu/'],
  [CT_EVENT_PART, '/events/cu/'],
  [CT_FULL_LESSON, '/events/cu/'],
  [CT_FRIENDS_GATHERING, '/events/cu/'],
  [CT_MEAL, '/events/cu/'],
  [CT_ARTICLE, '/publications/articles/cu/'],
]);

/* WARNING!!!
   This function MUST be synchronized with the next one: canonicalContentType
 */
export const canonicalLink = (entity, mediaLang) => {
  const base = 'https://kabbalahmedia.info';
  if (!entity) {
    return base + '/';
  }

  // source
  if (entity.content_type === 'SOURCE') {
    return base + `/sources/${entity.id}`;
  }

  if (entity.content_type === 'POST') {
    const [blogID, postID] = entity.id.split('-');
    const blogName         = blogNames.get(parseInt(blogID, 10)) || 'laitman-co-il';

    return base + `/publications/blog/${blogName}/${postID}`;
  }

  // collections
  switch (entity.content_type) {
  case CT_DAILY_LESSON:
  case CT_SPECIAL_LESSON:
    return base + `/lessons/daily/c/${entity.id}`;
  case CT_VIRTUAL_LESSONS:
    return base + `/lessons/virtual/c/${entity.id}`;
  case CT_LECTURE_SERIES:
    return base + `/lessons/lectures/c/${entity.id}`;
  case CT_WOMEN_LESSONS:
    return base + `/lessons/women/c/${entity.id}`;
    // case CT_CHILDREN_LESSONS:
    //   return base + `/lessons/children/c/${entity.id}`;
  case CT_LESSONS_SERIES:
    return base + `/lessons/series/c/${entity.id}`;
  case CT_VIDEO_PROGRAM:
  case CT_CLIPS:
    return base + `/programs/c/${entity.id}`;
  case CT_ARTICLES:
    return base + `/publications/articles/c/${entity.id}`;
  case CT_FRIENDS_GATHERINGS:
  case CT_MEALS:
  case CT_CONGRESS:
  case CT_HOLIDAY:
  case CT_PICNIC:
  case CT_UNITY_DAY:
    return base + `/events/c/${entity.id}`;
  default:
    break;
  }

  // units whose canonical collection is an event goes as an event item
  const collection = canonicalCollection(entity);
  if (collection && EVENT_TYPES.indexOf(collection.content_type) !== -1) {
    return base + `/events/cu/${entity.id}`;
  }

  const mediaLangSuffix = mediaLang ? `?language=${mediaLang}` : '';

  // unit based on type
  const prefix = mediaPrefix.get(entity.content_type);
  if (prefix) {
    return base + `${prefix}${entity.id}${mediaLangSuffix}`;
  } else {
    return base + '/';
  }
};

/* WARNING!!!
   This function MUST be synchronized with the previous one: canonicalLink
 */
export const canonicalContentType = (entity) => {
  switch (entity) {
  case 'sources':
    return ['SOURCE'];
  case 'lessons':
    return [...COLLECTION_LESSONS_TYPE, ...UNIT_LESSONS_TYPE];
  case 'programs':
    return [...COLLECTION_PROGRAMS_TYPE, ...UNIT_PROGRAMS_TYPE];
  case 'publications':
    return ['POST', CT_ARTICLES, ...COLLECTION_PUBLICATIONS_TYPE, ...UNIT_PUBLICATIONS_TYPE];
  case 'events':
    return [...COLLECTION_EVENTS_TYPE, ...UNIT_EVENTS_TYPE];
  default:
    return [];
  }
};
