//region **** Common classes ****
class fw2fn {

    static init() {
        const _this = this
    }

    static id_to_int(str) { // covert str "aaa123" or int to int 123
        if (Number.isInteger(str)) {
            return str
        }
        let numb = str.match(/(\d)+/g)
        try {
            return parseInt(numb[numb.length - 1], 10);
        } catch {
            return null;
        }
    }

    static id_to_arr_int(d) { // covert str or [str1, str2] to int array
        if (!Array.isArray(d)) {
            d = [d]
        }
        let ret = []
        for (let item of d) {
            let int_item = this.id_to_int(item)
            if (int_item !== null) {
                ret.push(int_item)
            }
        }
        return ret
    }

    static disable_elements(disable = true) {
        let controls = document.querySelectorAll("button, input, select, textarea");
        for (let c of controls) {
            c.disabled = disable;
        }
    }
}


class FW2 {
    constructor(url, token, keepalive = 0, wss = false) {
        const _this = this
        this._url = url
        this._keepalive(keepalive)
        this._wss = wss
        // RPC
        this.notify = new EventTarget();
        this._rpc_id = 0
        this._receive_rpc_fn_dict = {}  // {rpc_id: fn}
        // Global vars
        this._ws_id = null
        this._token = token
        this._connect_f = false
        // Public functions
        this.get_table = async function (...args) {
            return await _this._get_table(...args)
        }
        // on close browser
        window.onbeforeunload = function () {
            _this._ws.send('close');
        }
        this._ping(10)
    }

    async init() {
        const _this = this
        let protocol = 'ws'
        if (this._wss) protocol = 'wss'
        this._ws = new WebSocket(`${protocol}://${location.hostname}:${location.port}/${this._url}`)
        this._ws.onopen = () => {
            this._connect_f = true
            this._ws.send(this.token)
        }
        this._ws.onmessage = (m) => {
            this._receive_rpc(m)
        }

        this._ws.onclose = () => {
            this._on_close_ws()
        }

    }

    _receive_rpc(m) {
        let msg = JSON.parse(m.data)
        if (msg.jsonrpc !== '2.0') return false
        if (msg.method) {  // this is Notify
            this.notify.dispatchEvent(new CustomEvent(msg.method, {detail: msg.params}))
        } else {
            if (msg.id > 0) {
                let err = false
                if (msg.error) {
                    err = msg.error
                }
                this._receive_rpc_fn_dict[msg.id](msg.result, err) // todo send error
            }
        }
    }

    call_rpc(rpc_name, params = {}, error = true) {
        const _this = this
        return new Promise((resolve, reject) => {
            if (!this._connect_f) {
                if (error) {
                    reject('ws not connected')
                } else {
                    resolve(false)
                }
            }
            this._rpc_id++;
            const rpc_id = this._rpc_id
            let rpc = JSON.stringify({jsonrpc: "2.0", method: rpc_name, "params": params, "id": rpc_id})
            this._ws.send(rpc)
            this._receive_rpc_fn_dict[rpc_id] = (d, err) => {
                if (!err) {
                    resolve(d)
                } else {
                    if (error) {
                        reject(err)
                    } else {
                        resolve(false)
                    }
                }
                delete _this._receive_rpc_fn_dict[rpc_id]
            }
        })
    }

    get ws_id() {
        return this._ws_id
    }

    get token() {
        return this._token
    }

    _keepalive(interval) {
        if (interval > 0) {
            setInterval(async function f() {
                await fetch('/keepalive')

            }, interval * 1000)
        }
    }

    _ping(interval) {
        const _this = this

        if (interval > 0) {
            setInterval(async function f() {
                if (_this._connect_f) {
                    _this._ws.send('')
                }
            }, interval * 1000)
        }
    }

    _on_close_ws() {
        this._connect_f = false
        document.body.innerHTML = `<div class="row">
    <div class="col-sm-4"></div>
    <div class="col-sm-4">
        <div class="card text-center">
            <span style="width: 6rem" class="card-img-top m-2">
                <svg viewBox="0 0 512 512" xml:space="preserve" xmlns="http://www.w3.org/2000/svg"><path id="a" d="m463 319c14.4 23.2 15.2 52.4 1.92 76.2-13.2 23.8-38.4 38.7-65.7 38.7h-271c-27.3 0-52.4-14.8-65.7-38.7-13.2-23.8-12.5-53 1.92-76.2l136-218c13.8-22.2 37.7-35.4 63.8-35.4s50 13.2 63.8 35.4z" fill="#ffd400"/><path id="b" d="m425 342c12.5 20.1-1.93 46.1-25.6 46.1h-271c-23.6 0-38.1-26-25.6-46.1l136-218c11.8-18.9 39.4-18.9 51.2 0z" fill="#3f4751"/><path id="c" d="m278 324c3.62 3.62 5.86 8.62 5.86 14.1s-2.24 10.5-5.86 14.1-8.62 5.86-14.1 5.86c-11 0-20-8.95-20-20 0-5.52 2.24-10.5 5.86-14.1s8.62-5.86 14.1-5.86 10.5 2.24 14.1 5.86z" fill="#ffd400"/><path id="d" d="m284 188v85c0 11-8.95 20-20 20s-20-8.95-20-20v-85c0-11 8.95-20 20-20 5.52 0 10.5 2.24 14.1 5.86s5.86 8.62 5.86 14.1z" fill="#ffd400"/></svg>
            </span>
            <div class="card-body">
                <h5 class="card-title">Сессия закрыта</h5>
                <a href="#" class="btn btn-primary" onclick="location.reload()">Новая сессия</a>
            </div>
        </div>
    </div>
    <div class="col-sm-4"></div>
</div>`
    }

    async _get_table(url, config, params) {
        let message = await this.call_rpc(url, params)
        if (params['field_arr']) {
            let message_named = []
            message.forEach((element) => {
                let record = {}
                params['field_arr'].forEach((field_name, i) => {
                    record[field_name] = element[i]
                })
                message_named.push(record)
            })
            message = message_named
        }
        if (params['format_fn']) {
            return params['format_fn'](message)
        }
        return message
    }


}

class StoredInputs {

    constructor(fw, rpc_name) {
        const _this = this
        this._fw = fw
        this.rpc_name = rpc_name
        this.input_dict = new Map()

    }

    append(el) {
        el.rpc_name = this.rpc_name
        el.addEventListener('input', this.on_input)
        el.addEventListener('change', e => this.on_change(e))
        el.addEventListener('focus', this.on_focus)
        el.addEventListener('focusout', this.on_focusout)
        this.input_dict.set(el.id, el)
        new bootstrap.Popover(el, {trigger: 'manual'})
    }

    load() {
        const _this = this
        this._fw.call_rpc(this.rpc_name, {inputs: Array.from(this.input_dict.keys())})
            .then(function (message) {
                // _this.input.value = message
                for (let [key, obj] of _this.input_dict) {
                    obj.value = message[key]
                }
            }, function (err) {

            });

    }

    on_focus(e) {
        let popover = bootstrap.Popover.getInstance(e.target)
        if (e.target.err_value) {
            popover.show()
        }
    }

    on_focusout(e) {
        let popover = bootstrap.Popover.getInstance(e.target)
        if (popover._popper) {
            popover.hide()
        }
    }

    on_input(e) {
        let svg = `data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='12' height='12'><polygon points='0,0 12,0 12,12' style='fill:rgba(180,0,0,0.7)' /></svg>`
        e.target.style.background = `url("${svg}") no-repeat top right`
    }

    async on_change(e) {
        const _this = this

        try {
            let message = await fw.call_rpc(e.target.rpc_name, {value: e.target.value, id: e.target.id})
            e.target.style.background = ''
            e.target.err_value = false
        } catch (err) {
            e.target.title = 'Ошибка'
            e.target.dataset.content = err
            e.target.err_value = true

        }

    }
}

class PhoneButton {
    // {asterisk_name: {phone_number: [PhoneButton, PhoneButton, ...]}}
    static _phone_numbers
    static _first_init = true


    constructor(phone_number, asterisk, big = false) {
        if (!(asterisk in PhoneButton._phone_numbers)) {
            PhoneButton._phone_numbers[asterisk] = {}
        }
        if (!(phone_number in PhoneButton._phone_numbers[asterisk])) {
            PhoneButton._phone_numbers[asterisk][phone_number] = []
        }
        PhoneButton._phone_numbers[asterisk][phone_number].push(this)
        this._phone_number = phone_number
        this._asterisk = asterisk
        // init element
        this._el = this._init_element()
        this._change_status('broken')
        this.el.title = `Станция: ${this._asterisk}`
        if (PhoneButton._first_init) {
            PhoneButton._first_init = false
            PhoneButton._init_handler()
        }
    }

    _init_element() {
        let el = document.createElement('div')
        el.className = 'phone-button'
        el.dataset.phone_number = this._phone_number
        el.dataset.asterisk = this._asterisk
        el.textContent = this._phone_number
        return el
    }

    static _init_handler() {
        // Обработчик события изменения статуса
        fw_gui.notify.addEventListener('on_change_extension_status', (e) => {
            let params = e.detail
            console.log('on_change_extension_status:', params)
            // params = {state: [[asterisk_name, ext_name, device_state, traffic_state, simple_state], ...]}
            for (const state of params.state) {

                try {
                    for (const pb of PhoneButton._phone_numbers[state[0]][state[1]]) {
                        console.log('state:', state[4])
                        switch (state[4]) {
                            case -100:
                                pb._change_status('broken')
                                break
                            case 0:
                                pb._change_status('free')
                                break
                            case 5:
                                pb._change_status('ringing')
                                break
                            case 10:
                                pb._change_status('busy')
                                break
                        }
                    }
                } catch (e) {

                }
            }
        })
    }

    _change_status(status) {
        const status_arr = ['broken', 'free', 'ringing', 'busy']
        for (const status of status_arr) {
            this._el.classList.remove('phone-button-' + status)
        }
        this._el.classList.add('phone-button-' + status)
    }

    get el() {
        return this._el
    }

    get phone_number() {
        return this._phone_number
    }

    set change_status(arr) {

    }

    static subscribe() {
        // {asterisk: {phone_number1, phone_number2, ...}}}
        let ret_1 = {}

        for (const [asterisk, phone_number] of Object.entries(PhoneButton._phone_numbers)) {
            if (!(asterisk in ret_1)) {
                ret_1[asterisk] = new Set()
                for (const phone in phone_number) {
                    ret_1[asterisk].add(phone)
                }
            }
        }
        let ret = {}
        for (const [asterisk, phone_number_set] of Object.entries(ret_1)) {
            ret[asterisk] = Array.from(phone_number_set)
        }

        fw_gui.call_rpc('subscribe_all_phone_button', ret).then()
    }

}


class BootstrapDialogsLng {
    constructor(lng) {
        this._lng = lng
    }

    ok() {
        let w = {
            'en': 'Ok',
            'ru': 'Ok'
        }
        return this._ret(w)
    }

    cancel() {
        let w = {
            'en': 'Cancel',
            'ru': 'Отмена'
        }
        return this._ret(w)
    }

    _ret(w) {
        if (w[this._lng]) {
            return w[this._lng]
        } else {
            return w['en']
        }
    }

}

class BootstrapDialogs {
    constructor() {
        this._recalc_z_index()
        this._lng = new BootstrapDialogsLng('ru')
    }

    async confirm(title, body, params = {}) {
        if (body instanceof Element && params['onsubmit_fn']) {
            await this._confirm_form(title, body, params)
        } else {
            return await this._confirm_text(title, body, params)
        }
    }

    _show_dialog(modal) {
        modal.show()
    }

    async _confirm_text(title, body, params) {
        let only_ok = false
        if (params['only_ok']) {
            only_ok = true
        }
        let dlg = this._create_confirm_dialog(title, body, params, only_ok)
        let ok_button = dlg.btn_ok
        let dialog = new bootstrap.Modal(dlg.modal, {backdrop: 'static'})
        setTimeout(this._show_dialog, 0, dialog)

        return new Promise((resolve, reject) => {
            ok_button.onclick = function () {
                resolve(true)
                dialog.hide()
                dlg.modal.remove()
            }

            dlg.modal.addEventListener('hide.bs.modal', function (e) {
                resolve(false)
                dlg.modal.remove()
            })
        })
    }

    async _confirm_form(title, body, params) {
        const _this = this
        let dlg = this._create_confirm_dialog(title, body, params)
        let dialog = new bootstrap.Modal(dlg.modal, {backdrop: 'static'})
        let cln_frm = body.cloneNode(true);
        let submit_btn = document.createElement('button')
        submit_btn.type = 'submit'
        submit_btn.hidden = true

        cln_frm.onsubmit = async function (e) {
            e.preventDefault()
            let ser = _this._serialize_form(cln_frm)
            if (await params.onsubmit_fn(ser, cln_frm)) {
                dialog.hide()
                dlg.modal.remove()
            }
        }

        cln_frm.appendChild(submit_btn)
        cln_frm.hidden = false
        dlg.body_div.append(cln_frm)
        let ok_button = dlg.btn_ok
        setTimeout(this._show_dialog, 0, dialog)

        ok_button.onclick = function () {
            submit_btn.click()
        }

        dlg.modal.addEventListener('hide.bs.modal', function (e) {
            dlg.modal.remove()
            params.onsubmit_fn(false, false)
        })
    }

    _create_confirm_dialog(title, body, params, only_ok = false) {
        let modal_div = document.createElement('div')
        modal_div.className = 'modal fade'
        modal_div.dataset.backdrop = 'static'
        modal_div.role = 'dialog'
        modal_div.setAttribute('aria-hidden', 'true')
        let modal_dialog_div = document.createElement('div')
        modal_dialog_div.className = 'modal-dialog'
        let modal_content_div = document.createElement('div')
        modal_content_div.className = 'modal-content'
        let modal_header_div = document.createElement('div')
        modal_header_div.className = 'modal-header'
        modal_header_div.insertAdjacentHTML('beforeend', `<h5 class="modal-title">${title}</h5>
        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>`)
        modal_content_div.appendChild(modal_header_div)

        let modal_body_div = document.createElement('div')
        modal_body_div.className = 'modal-body'
        if (typeof body === 'string') {
            modal_body_div.insertAdjacentHTML('beforeend', body)
            if (body.length === 0) {
                modal_header_div.classList.add('border-0')
                modal_body_div.classList.add('p-0')
            }
        }
        modal_content_div.appendChild(modal_body_div)

        let modal_footer_div = document.createElement('div')
        modal_footer_div.className = 'modal-footer'
        if (!only_ok) {
            let btn_cancel = document.createElement('button')
            btn_cancel.className = 'btn btn-secondary'
            btn_cancel.textContent = this._lng.cancel()
            btn_cancel.dataset.bsDismiss = "modal"
            btn_cancel.type = "button"
            modal_footer_div.appendChild(btn_cancel)
        }
        let btn_ok = document.createElement('button')
        btn_ok.className = 'btn btn-primary'
        // set caption
        if (params.ok_caption !== undefined) {
            btn_ok.textContent = params.ok_caption
        } else {
            btn_ok.textContent = this._lng.ok()
        }
        btn_ok.type = "button"
        // append
        modal_footer_div.appendChild(btn_ok)
        modal_content_div.appendChild(modal_footer_div)
        modal_dialog_div.appendChild(modal_content_div)
        modal_div.appendChild(modal_dialog_div)
        document.body.appendChild(modal_div)
        return {modal: modal_div, btn_ok: btn_ok, body_div: modal_body_div}
    }

    _recalc_z_index() {
        document.addEventListener('shown.bs.modal', function (e) {
            let el = e.target
            let all_modal = document.querySelectorAll('.modal')
            let zIndex = 1040
            all_modal.forEach(function (el) {
                if (getComputedStyle(el).display !== 'none')
                    zIndex += 10
            })
            el.style.zIndex = zIndex.toString()
            setTimeout(function () {
                //$('.modal-backdrop').not('.modal-stack').css('z-index', zIndex - 1).addClass('modal-stack');
                let modal_backdrop = document.querySelectorAll('.modal-backdrop')
                modal_backdrop.forEach(function (el) {
                    if (!el.classList.contains('modal-stack')) {
                        el.style.zIndex = (zIndex - 1).toString()
                        el.classList.add('modal-stack')
                    }
                })
            }, 0);
        })
    }

    _serialize_form(frm) {
        let ret_dict = {}
        const selectors = frm.querySelectorAll("input")
        selectors.forEach(function (selector) {
            if (selector.dataset.name) {
                if (selector.type === 'checkbox') {
                    if (selector.checked && selector.name) {
                        try {
                            ret_dict[selector.name].push(selector.dataset.name)
                        } catch {
                            ret_dict[selector.name] = []
                            ret_dict[selector.name].push(selector.dataset.name)
                        }
                    }
                } else {
                    ret_dict[selector.dataset.name] = selector.value
                }
            } else {
                if (selector.type === 'radio') {
                    if (selector.checked && selector.name) {
                        ret_dict[selector.name] = selector.value
                    }
                }
            }

        });
        return ret_dict
    }


}

// need lodash.js
class BsMultiSelect {
    static id_cnt = 0;

    /**
     * Создает в выбранном контейнере виджет BsMultiSelect
     * @param container DOM element
     */
    constructor(container) {
        const _this = this
        this._container = container
        //
        this._first_data_f = true
        this._first_data = {}
        this._data = {}
        this._prev_data = {}
        this.onclick = undefined
    }

    add(items) {
        // {'item1': true, 'item2': false}
        // check if exist
        for (let key in items) {
            if (key in this._data) {
                delete items[key]
            } else {
                this._data[key] = items[key]
            }
        }

        if (this._first_data_f) {
            this._first_data = items
            this._first_data_f = false
            this._data = Object.assign({}, items);
        }

        for (let item in items) {
            let div = this._create_item(item, items[item])
            this._container.appendChild(div)
        }
    }

    _create_item(item_name, val) {
        let item_id = BsMultiSelect.id_cnt++
        const _this = this
        let div = document.createElement('div')
        div.className = 'form-check'
        div.type = 'checkbox'
        div.id = `bs_multiselect_${item_id}`
        div.dataset.id = `${item_id}`
        let input = document.createElement('input')
        input.className = 'form-check-input'
        input.type = "checkbox"
        input.checked = val
        input.id = div.id + '_i'
        input.dataset.item_name = item_name
        let label = document.createElement('label')
        label.className = 'form-check-label'
        label.setAttribute('for', input.id)
        label.textContent = item_name
        div.appendChild(input)
        div.appendChild(label)
        // events
        input.onchange = function (e) {
            _this._item_on_change(e.target)
        }
        return div
    }

    async _item_on_change(el) {
        if (this.onclick !== undefined) {
            try {
                let ret = this.onclick([el.dataset.item_name, el.checked])
                if (ret === false) {
                    el.checked = !el.checked
                    return
                }
            } catch (e) {
                console.error(e)
            }
        }
        this._data[el.dataset.item_name] = el.checked
    }

    get data() {
        return this._data
    }

    get is_changed() {
        return !_.isEqual(this._first_data, this._data)
    }

    get checked_array() {
        let ret = []
        for (let key in this._data) {
            if (this._data[key]) {
                ret.push(key)
            }
        }
        return ret
    }

    get unchecked_array() {
        let ret = []
        for (let key in this._data) {
            if (!this._data[key]) {
                ret.push(key)
            }
        }
        return ret
    }

}


class BsMultiSelectWs {

}


class CtxMenu {
    static open_menu = null

    constructor(el, menu) {
        const _this = this
        this._el = el
        this.event = new EventTarget();
        this._menu = menu
        this._menu_dom = this._create_menu()

        this._menu_dom.onclick = (e) => {
            _this._click_menu(e, _this)
        }

        this._el.oncontextmenu = (e) => {
            try {
                CtxMenu.open_menu.remove()
            } catch (e) {
            }
            e.preventDefault()
            document.body.appendChild(_this._menu_dom)
            let top, left
            let sticky = _this._menu.sticky
            if (sticky) {
                let coord = _this._el.getBoundingClientRect()
                let coord_menu = _this._menu_dom.getBoundingClientRect()
                if (sticky === 'bottom') {
                    top = coord.y + coord.height + pageYOffset
                    left = coord.x + pageXOffset
                }
                if (sticky === 'top') {
                    top = coord.y - coord_menu.height + pageYOffset
                    left = coord.x + pageXOffset
                }

            } else {
                top = e.pageY - 10
                left = e.pageX + 10
            }
            _this._menu_dom.style.top = `${top}px`
            _this._menu_dom.style.left = `${left}px`
            _this._menu_dom.style.boxShadow = '0 0 5px rgba(0,0,0,0.3)'
            _this._menu_dom.style.display = 'block'
            _this._menu_dom.style.position = 'absolute'
            _this._menu_dom.style.zIndex = 1000
            document.addEventListener('click', function (e) {
                if (!_this._menu_dom.contains(e.target)) {
                    _this._menu_dom.remove()
                }
            })
            CtxMenu.open_menu = _this._menu_dom
        }
    }

    _click_menu(e, _this) {
        let el = e.target
        if (!el.classList.contains('list-group-item')) {
            el = el.closest('li');
        }
        _this._menu_dom.remove()
        document.removeEventListener('click', () => {
        })
        _this.event.dispatchEvent(new CustomEvent("onclick", {
            detail:
                {id: el.dataset.id, owner: _this._el},
        }))
        if (typeof _this._menu.callback === 'function') {
            _this._menu.callback(_this._el, el.dataset.id)
        }
    }

    _create_menu() {
        let cont = document.createElement('ul')
        cont.className = 'list-group'
        this._menu.menu.forEach((v) => {
            let li = document.createElement('li')
            li.className = 'list-group-item list-group-item-action'
            if (v.add_class) {
                v.add_class.forEach((cls) => {
                    li.classList.add(cls)
                })
            }
            if (v.id === 'separator') {
                li.style.padding = '0'
            } else {
                li.style.padding = '3px 10px 3px 10px'
                li.style.cursor = 'pointer'
                li.style.userSelect = 'none'
                li.innerHTML = `${v.text}`
                li.dataset.id = v.id
            }

            cont.appendChild(li)
        })
        return cont
    }
}

class BootstrapToast {
    static _fw
    static _toast_id_cnt = 0
    static _toast_id_BootstrapToast = {}
    static _toast_container_div
    static default_option = {'header_bg_color': ''};

    static init(fw, notify_name, zIndex = 2000) {
        const _this = this
        this._fw = fw
        this._toast_container_div = document.createElement('div')
        this._toast_container_div.className = 'toast-container position-absolute top-0 end-0 p-3'
        this._toast_container_div.style.zIndex = zIndex
        document.body.append(this._toast_container_div)
        this._fw.notify.addEventListener('bootstrap_toast', (e) => {
            let params = e.detail
            new this(params[0], params[1], params[2])
        })
    }

    constructor(header_html = '', body_html = '', option = {}) {
        this._option = Object.assign({}, BootstrapToast.default_option, option);
        BootstrapToast._toast_id_cnt++
        const _this = this
        this._id = BootstrapToast._toast_id_cnt
        this._el = document.createElement('div')
        this._el.className = 'toast'
        this._el.setAttribute('role', 'alert')
        this._el.setAttribute('aria-live', 'assertive')
        this._el.setAttribute('aria-atomic', 'true')
        this._create_toast(header_html, body_html, this._option)
        BootstrapToast._toast_container_div.appendChild(this._el)
        //
        this._bs = new bootstrap.Toast(this._el, this._option)
        BootstrapToast._toast_id_BootstrapToast[this._id] = this
        this._el.addEventListener('hidden.bs.toast', function () {
            console.log('hidden.bs.toast:', _this._id)
            _this._el.remove()
            try {
                delete BootstrapToast._toast_id_BootstrapToast[_this._id]
            } catch {

            }
        })
        this.show()
    }

    _bg_color_class(option) {
        if (option === 'danger') {
            return ['bg-danger', 'text-white']
        } else if (option === 'warning') {
            return ['bg-warning', 'text-dark']
        } else if (option === 'info') {
            return ['bg-info', 'text-dark']
        }
        return []
    }

    _create_toast(header_html = '', body_html = '', option) {
        let toast_header_div = document.createElement('div')
        toast_header_div.className = 'toast-header'
        toast_header_div.classList.add(...this._bg_color_class(option.header_bg_color))
        toast_header_div.insertAdjacentHTML("beforeend", `<strong class="me-auto">${header_html}</strong>
        <small class="text-muted">just now</small>
        <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>`)
        this._el.appendChild(toast_header_div)
        let toast_body_div = document.createElement('div')
        toast_body_div.className = 'toast-body'
        toast_body_div.innerHTML = body_html
        this._el.appendChild(toast_body_div)
    }

    show() {
        this._bs.show()
    }

}

class TabulatorLocale {
    static langs = {
        "ru-ru": {
            "columns": {
                "name": "Name", //replace the title of column name with the value "Name"
            },
            "ajax": {
                "loading": "Loading", //ajax loader text
                "error": "Error", //ajax error text
            },
            "groups": { //copy for the auto generated item count in group header
                "item": "item", //the singular  for item
                "items": "items", //the plural for items
            },
            "pagination": {
                "page_size": "Page Size", //label for the page size select element
                "page_title": "Show Page",//tooltip text for the numeric page button, appears in front of the page number (eg. "Show Page" will result in a tool tip of "Show Page 1" on the page 1 button)
                "first": "First", //text for the first page button
                "first_title": "First Page", //tooltip text for the first page button
                "last": "Last",
                "last_title": "Last Page",
                "prev": "Prev",
                "prev_title": "Prev Page",
                "next": "Next",
                "next_title": "Next Page",
                "all": "All",
            },
            "headerFilters": {
                "default": "поиск...", //default header filter placeholder text
                "columns": {
                    "name": "filter name...", //replace default header filter text for column name
                }
            }
        }
    }
}

class PostJsonRpc {
    constructor(url, cookies = true) {
        const _this = this
        this._url = url
        this._cookies = cookies
        this._i = 0

    }

    async send(rpc_name, data) {
        let ret
        this._i++
        let credentials = "omit"
        if (this._cookies) {
            credentials = "same-origin"
        }
        try {
            let d = await fetch(`${this._url}`, {
                method: 'POST',
                credentials: credentials,
                headers: {
                    'Content-type': 'application/text; charset=UTF-8' // The type of data you're sending
                },
                body: `{"jsonrpc": "2.0", "method": "${rpc_name}", "params": ${JSON.stringify(data)}, "id": ${this._i}}`
            })
            ret = await d.json()
            ret = ret.result
        } catch (error) {
            return false
        }
        return ret
    }

}


//endregion