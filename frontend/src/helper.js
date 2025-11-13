import chaturbate from './assets/images/logo_chaturbate.png';
import fansly from './assets/images/logo_fansly.png';
import stripchat from './assets/images/logo_stripchat.png';

export const eventTypes = [
  'tip', 'follow', 'unfollow', 'subscribe'
]

export const imageLookup = {
  chaturbate,
  fansly,
  stripchat,
}

export const genderLookup = {
  "m": "Male",
  "f": "Female",
  "t": "Trans",
  "c": "Couple",
  "u": "Unknown",
}

export const getSubscriptionName = (tier) => {
  return {
    chaturbate: {
      name: "fanclub",
      shortText: "joined fanclub",
      longText: "has joined your fanclub",
    },
    fansly: {
      name: tier,
      shortText: "subscribed to " + tier,
      longText: `has subscribed to Tier "${tier}"`,
    },
    fansly: {
      name: tier,
      shortText: "join fanclub tier " + tier,
      longText: `has joined your Fanclub, Tier "${tier}"`,
    },
  }
}

export function getEventItemText(event) {
  switch (event.Type) {
    case "tip":
      return {
        shortText: event.Event.TipValue + event.Event.TipCurrencySymbol,
        longText: `User ${event.Event.User.Username} has tipped ${event.Event.TipValue + " " + event.Event.TipCurrency}`,
      }
    case "follow":
      return {
        shortText: "followed",
        longText: `User ${event.Event.User.Username} has followed you`,
      }
    case "unfollow":
      return {
        shortText: "unfollowed",
        longText: `User ${event.Event.User.Username} has unfollowed you`,
      }
    case "subscribe":
      const subTexts = getSubscriptionName(event.Event.User.Tier)[event.Provider];
      return {
        shortText: subTexts.shortText,
        longText: `User ${event.Event.User.Username} has ${subTexts.longText}`,
      }
  }
}
export function capitalizeFirstLetter(string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

