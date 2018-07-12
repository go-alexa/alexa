// Package response allows constructing responses to Alexa requests.
package response

import (
	"github.com/go-alexa/alexa/parser"
)

const (
	// OutputSpeechPlain means the output will be plain text.
	OutputSpeechPlain = "PlainText"
	// OutputSpeechSSML means the output will be SSML.
	OutputSpeechSSML = "SSML"

	// CardSimple means a simple card (only title and body).
	CardSimple = "Simple"
	// CardStandard means a standard card (title, body, and image).
	CardStandard = "Standard"
	// CardLinkAccount is a card for linking the user's account.
	CardLinkAccount = "LinkAccount"
	// DialogDelegateDirective is the dialog delegate directive type
	DialogDelegateDirective = "Dialog.Delegate"
	// ElicitSlotDirective is the elicit slot directive type
	ElicitSlotDirective = "Dialog.ElicitSlot"
	// ConfirmSlotDirective is the confirm slot directive type
	ConfirmSlotDirective = "Dialog.ConfirmSlot"
	// ConfirmIntentDirective is the confirm intent directive type
	ConfirmIntentDirective = "Dialog.ConfirmIntent"
)

// Response is the base response struct.
type Response struct {
	Version    string                   `json:"version"`
	Attributes parser.SessionAttributes `json:"sessionAttributes"`
	Response   InnerResponse            `json:"response"`
}

// InnerResponse is all the actual information for the response.
type InnerResponse struct {
	OutputSpeech     *OutputSpeech `json:"outputSpeech,omitempty"`
	Card             *Card         `json:"card,omitempty"`
	Reprompt         *Reprompt     `json:"reprompt,omitempty"`
	Directives       []*Directive  `json:"directives,omitempty"`
	ShouldEndSession bool          `json:"shouldEndSession"`
}

// OutputSpeech is actual spoken text for the response.
type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	SSML string `json:"ssml,omitempty"`
}

// Card is any information needed to return a card. Not all fields are required
// for each type of card.
type Card struct {
	Type      string     `json:"type"`
	Title     string     `json:"title,omitempty"`
	Content   string     `json:"content,omitempty"`
	Text      string     `json:"text,omitempty"`
	ImageURLs *ImageURLs `json:"image,omitempty"`
}

// ImageURLs are URLs for images for the standard card.
type ImageURLs struct {
	SmallImageURL string `json:"smallImageUrl"`
	LargeImageURL string `json:"largeImageUrl"`
}

// Reprompt is the speech for a reprompt message, if there is one.
type Reprompt struct {
	OutputSpeech OutputSpeech `json:"outputSpeech,omitempty"`
}

// Directive is a Dialog Directive
// Ref: https://developer.amazon.com/docs/custom-skills/dialog-interface-reference.html
type Directive struct {
	Type          string         `json:"type"`
	Intent        *parser.Intent `json:"updatedIntent,omitempty"`
	SlotToElicit  string         `json:"slotToElicit,omitempty"`
	SlotToConfirm string         `json:"slotToConfirm,omitempty"`
}

// New creates a new Response with some default values set.
func New() *Response {
	return &Response{
		Version: "1.0",
		Response: InnerResponse{
			ShouldEndSession: true,
			OutputSpeech:     nil,
		},
	}
}

// AddSpeech adds a simple text response.
func (r *Response) AddSpeech(speech string) *Response {
	r.Response.OutputSpeech = &OutputSpeech{
		Type: OutputSpeechPlain,
		Text: speech,
	}

	return r
}

// AddSSMLSpeech adds a SSML text response.
func (r *Response) AddSSMLSpeech(speech string) *Response {
	r.Response.OutputSpeech = &OutputSpeech{
		Type: OutputSpeechSSML,
		SSML: speech,
	}

	return r
}

// AddCard adds a simple card response.
func (r *Response) AddCard(title, content string) *Response {
	r.Response.Card = &Card{
		Type:    CardSimple,
		Title:   title,
		Content: content,
	}

	return r
}

// AddStandardCard adds a standard card response.
func (r *Response) AddStandardCard(title, text, smallImageURL, largeImageURL string) *Response {
	r.Response.Card = &Card{
		Type:  CardStandard,
		Title: title,
		Text:  text,
		ImageURLs: &ImageURLs{
			SmallImageURL: smallImageURL,
			LargeImageURL: largeImageURL,
		},
	}

	return r
}

// AddLinkAccountCard adds a link account card.
func (r *Response) AddLinkAccountCard() *Response {
	r.Response.Card = &Card{
		Type: CardLinkAccount,
	}

	return r
}

// AddReprompt adds a plain text reprompt message.
func (r *Response) AddReprompt(speech string) *Response {
	r.Response.Reprompt = &Reprompt{
		OutputSpeech: OutputSpeech{
			Type: OutputSpeechPlain,
			Text: speech,
		},
	}

	return r
}

// AddSSMLReprompt adds a SSML reprompt message.
func (r *Response) AddSSMLReprompt(speech string) *Response {
	r.Response.Reprompt = &Reprompt{
		OutputSpeech: OutputSpeech{
			Type: OutputSpeechSSML,
			Text: speech,
		},
	}

	return r
}

// AddDialogDelegateDirective adds a dialog delegate directive with the specified updated intent (If provided)
func (r *Response) AddDialogDelegateDirective(updatedIntent *parser.Intent) *Response {
	return r.addDirective(&Directive{
		Type:   DialogDelegateDirective,
		Intent: updatedIntent,
	})
}

// AddElicitSlotDirective adds a elicit slot directive with the slot to elicit and the updated intent (If provided)
func (r *Response) AddElicitSlotDirective(slotToElicit string, updatedIntent *parser.Intent) *Response {
	return r.addDirective(&Directive{
		Type:         ElicitSlotDirective,
		Intent:       updatedIntent,
		SlotToElicit: slotToElicit,
	})
}

// AddConfirmSlotDirective adds a confirm slot directive with the slot to confirm and the updated intent (If provided)
func (r *Response) AddConfirmSlotDirective(updatedIntent *parser.Intent, slotToConfirm string) *Response {
	return r.addDirective(&Directive{
		Type:          ConfirmSlotDirective,
		Intent:        updatedIntent,
		SlotToConfirm: slotToConfirm,
	})
}

// AddConfirmIntentDirective adds a confirm intent directive with the specified updated intent (If provided)
func (r *Response) AddConfirmIntentDirective(updatedIntent *parser.Intent) *Response {
	return r.addDirective(&Directive{
		Type:   ConfirmIntentDirective,
		Intent: updatedIntent,
	})
}

// SetAttributes sets attributes for the Session data. Note that this must be
// set with all data for each request, it does not merge or save data.
func (r *Response) SetAttributes(attrs parser.SessionAttributes) *Response {
	r.Attributes = attrs

	return r
}

// KeepAlive keeps a session alive instead of ending it.
func (r *Response) KeepAlive() *Response {
	r.Response.ShouldEndSession = false

	return r
}

func (r *Response) addDirective(directive *Directive) *Response {
	r.Response.Directives = append(r.Response.Directives, directive)
	return r
}
