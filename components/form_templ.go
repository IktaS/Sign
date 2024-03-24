// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.639
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Form() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form action=\"/sign\" method=\"post\" enctype=\"multipart/form-data\"><label for=\"username\">Username</label> <input type=\"username\" id=\"username\" name=\"username\" placeholder=\"Username of signer\"> <small>Ask IktaS for the username specific to you</small> <label for=\"password\">Password</label> <input type=\"password\" id=\"password\" name=\"password\" placeholder=\"Password of signer\"> <small>Ask IktaS for the password specific to you</small> <label for=\"file\">File to sign <input type=\"file\" id=\"file\" name=\"file\" accept=\".pdf, application/pdf\" onchange=\"onFileChange(this, event)\"></label><figure hidden id=\"pdf-editor\"><label for=\"qr-location-width\">Location Width <input type=\"range\" min=\"0\" max=\"100\" value=\"80\" id=\"qr-location-width\" name=\"qr-location-width\" onchange=\"reRenderPDF()\"></label> <label for=\"qr-location-height\">Location height <input type=\"range\" min=\"0\" max=\"100\" value=\"10\" id=\"qr-location-height\" name=\"qr-location-height\" onchange=\"reRenderPDF()\"></label><div class=\"grid\"><label for=\"qr-size\">Size (in px) <input type=\"text\" id=\"qr-size\" name=\"qr-size\" placeholder=\"100\" value=\"100\" onchange=\"reRenderPDF()\"></label> <label for=\"all-page\"><input type=\"checkbox\" id=\"all-page\" name=\"all-page\" checked onchange=\"reRenderPDF()\"> All Page</label> <label hidden for=\"qr-page\" id=\"qr-page-container\">Page <input type=\"text\" id=\"qr-page\" name=\"qr-page\" placeholder=\"1\" value=\"1\" onchange=\"reRenderPDF()\"></label></div><legend><strong>Preview</strong></legend> <iframe id=\"pdf\" style=\"margin: 1vh; width: 100%; height: 50vh;\"></iframe></figure><input type=\"submit\" value=\"Submit\"></form>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}