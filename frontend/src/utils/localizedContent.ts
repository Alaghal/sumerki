import type { TFunction } from 'i18next';

import type { EventChoice, KingdomEvent, MissionReport } from '../api/client';

function text(value: unknown): string {
  return typeof value === 'string' ? value : '';
}

export function getLocalizedEventTitle(t: TFunction, event: KingdomEvent) {
  return text(t(`events:content.${event.eventKey}.title`, { defaultValue: event.title || t('events:unknown.title') }));
}

export function getLocalizedEventBody(t: TFunction, event: KingdomEvent) {
  return text(t(`events:content.${event.eventKey}.body`, { defaultValue: event.body || t('events:unknown.body') }));
}

export function getLocalizedEventChoiceLabel(t: TFunction, event: KingdomEvent, choice: EventChoice) {
  return text(t(`events:content.${event.eventKey}.choices.${choice.key}.label`, { defaultValue: choice.label || t('events:unknown.choice') }));
}

export function getLocalizedEventChoiceDescription(t: TFunction, event: KingdomEvent, choice: EventChoice) {
  return text(t(`events:content.${event.eventKey}.choices.${choice.key}.description`, { defaultValue: choice.description || t('events:unknown.choiceDescription') }));
}

export function getLocalizedEventResultTitle(t: TFunction, event: KingdomEvent) {
  const fallback = event.result?.title || event.title || t('events:unknown.resultTitle');
  if (!event.selectedChoiceKey) {
    return text(fallback);
  }
  return text(t(`events:content.${event.eventKey}.choices.${event.selectedChoiceKey}.result.title`, { defaultValue: fallback }));
}

export function getLocalizedEventResultBody(t: TFunction, event: KingdomEvent) {
  const fallback = event.result?.body || event.body || t('events:unknown.resultBody');
  if (!event.selectedChoiceKey) {
    return text(fallback);
  }
  return text(t(`events:content.${event.eventKey}.choices.${event.selectedChoiceKey}.result.body`, { defaultValue: fallback }));
}

export function getLocalizedEventSelectedChoiceLabel(t: TFunction, event: KingdomEvent) {
  const choice = event.choices.find((candidate) => candidate.key === event.selectedChoiceKey);
  if (choice) {
    return getLocalizedEventChoiceLabel(t, event, choice);
  }
  if (!event.selectedChoiceKey) {
    return text(t('common:states.unknown'));
  }
  return text(t(`events:content.${event.eventKey}.choices.${event.selectedChoiceKey}.label`, { defaultValue: t('common:states.unknown') }));
}

export function getLocalizedReportTitle(t: TFunction, report: MissionReport) {
  return text(t(`reports:templates.${report.type}.title.${report.result}`, { defaultValue: report.title || t('reports:templates.fallback.title') }));
}

export function getLocalizedReportPhaseTitle(t: TFunction, phaseTitle: string) {
  return text(t(`reports:phaseTitles.${phaseTitle}`, { defaultValue: phaseTitle || t('reports:phaseTitles.unknown') }));
}
